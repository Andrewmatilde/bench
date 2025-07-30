package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"bench-server/pkg/config"
	"bench-server/pkg/database"
)

type Server struct {
	db     *sql.DB
	mux    *http.ServeMux
	config *config.Config
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

	server := &Server{
		db:     db,
		mux:    http.NewServeMux(),
		config: cfg,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	// 健康检查
	s.mux.HandleFunc("/health", s.healthHandler)

	// 传感器数据路由
	s.mux.HandleFunc("/api/sensor-data", s.sensorDataHandler)
	s.mux.HandleFunc("/api/stats", s.statsHandler)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

func (s *Server) sensorDataHandler(w http.ResponseWriter, r *http.Request) {
	return
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	if r.Method != "POST" {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data database.SensorData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// 数据验证
	if data.DeviceID == "" || data.MetricName == "" || data.Timestamp == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 验证优先级
	if data.Priority < 1 || data.Priority > 3 {
		data.Priority = 2 // 默认中等优先级
	}

	dbService := database.NewService(s.db)
	if err := dbService.InsertSensorData(&data); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	return
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	// 恢复 panic
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	if r.Method != "GET" {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dbService := database.NewService(s.db)
	stats, err := dbService.GetStats()
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.mux,
		ReadTimeout:  parseDuration(s.config.ReadTimeout),
		WriteTimeout: parseDuration(s.config.WriteTimeout),
		IdleTimeout:  parseDuration(s.config.IdleTimeout),
	}

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Printf("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server on port %s", s.config.Port)
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
