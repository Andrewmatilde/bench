package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"bench-server/pkg/config"
	"bench-server/pkg/database"
	"bench-server/pkg/handlers"
)

type Server struct {
	db      *sql.DB
	router  *mux.Router
	logger  *logrus.Logger
	config  *config.Config
	handler *handlers.Server
}

func NewServer(cfg *config.Config) (*Server, error) {
	// 初始化数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 配置数据库连接池
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 初始化表结构
	if err := database.InitDatabase(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 初始化日志
	logger := logrus.New()

	// 根据配置设置日志格式
	if cfg.LogFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	// 根据配置设置日志级别
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 创建处理器实例
	handlerServer := handlers.NewServer(db, logger)

	server := &Server{
		db:      db,
		router:  mux.NewRouter(),
		logger:  logger,
		config:  cfg,
		handler: handlerServer,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	// 健康检查
	s.router.HandleFunc("/health", s.healthHandler).Methods("GET")

	// 传感器数据路由
	s.router.HandleFunc("/api/sensor-data", s.handler.SensorDataHandler).Methods("POST")
	s.router.HandleFunc("/api/sensor-rw", s.handler.SensorReadWriteHandler).Methods("POST")
	s.router.HandleFunc("/api/batch-sensor-rw", s.handler.BatchSensorReadWriteHandler).Methods("POST")
	s.router.HandleFunc("/api/stats", s.handler.StatsHandler).Methods("GET")
	s.router.HandleFunc("/api/get-sensor-data", s.handler.GetSensorDataHandler).Methods("POST")

	// 添加中间件
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.recoveryMiddleware)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		s.logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"duration":   duration,
			"user_agent": r.UserAgent(),
		}).Info("HTTP request")
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.WithField("error", err).Error("Panic recovered")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(); err != nil {
		http.Error(w, "Database connection failed", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.router,
		ReadTimeout:  parseDuration(s.config.ReadTimeout),
		WriteTimeout: parseDuration(s.config.WriteTimeout),
		IdleTimeout:  parseDuration(s.config.IdleTimeout),
	}

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		s.logger.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			s.logger.WithError(err).Error("Server shutdown error")
		}
	}()

	s.logger.WithField("port", s.config.Port).Info("Starting server")
	return srv.ListenAndServe()
}

func (s *Server) Close() error {
	return s.db.Close()
}

func main() {
	cfg := config.New()

	server, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("Failed to parse duration %s: %v", s, err)
	}
	return d
}
