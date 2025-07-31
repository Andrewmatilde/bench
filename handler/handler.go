package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"bench-server/pkg/database"

	"github.com/sirupsen/logrus"
)

// Server HTTP服务器结构
type Server struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewServer 创建新的服务器实例
func NewServer(db *sql.DB, logger *logrus.Logger) *Server {
	return &Server{
		db:     db,
		logger: logger,
	}
}

// SensorDataHandler 处理传感器数据上报（扩展功能）
func (s *Server) SensorDataHandler(w http.ResponseWriter, r *http.Request) {
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
		s.logger.WithError(err).Error("Failed to insert sensor data")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Data inserted successfully",
	})
}

// StatsHandler 处理统计信息请求
func (s *Server) StatsHandler(w http.ResponseWriter, r *http.Request) {
	dbService := database.NewService(s.db)
	stats, err := dbService.GetStats()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get stats")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
