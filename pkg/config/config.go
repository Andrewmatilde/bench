package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// ConfigFile 配置文件结构
type ConfigFile struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Name         string `yaml:"name"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
	App struct {
		ReadTimeout    string `yaml:"read_timeout"`
		WriteTimeout   string `yaml:"write_timeout"`
		IdleTimeout    string `yaml:"idle_timeout"`
		BatchSize      int    `yaml:"batch_size"`
		FlushTimeout   string `yaml:"flush_timeout"`
		ChannelBuffer  int    `yaml:"channel_buffer"`
	} `yaml:"app"`
}

// Config 应用配置结构
type Config struct {
	Port         string `yaml:"port"`
	DBHost       string `yaml:"db_host"`
	DBPort       string `yaml:"db_port"`
	DBUser       string `yaml:"db_user"`
	DBPassword   string `yaml:"db_password"`
	DBName       string `yaml:"db_name"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	LogLevel     string `yaml:"log_level"`
	LogFormat    string `yaml:"log_format"`
	ReadTimeout   string `yaml:"read_timeout"`
	WriteTimeout  string `yaml:"write_timeout"`
	IdleTimeout   string `yaml:"idle_timeout"`
	BatchSize     int    `yaml:"batch_size"`
	FlushTimeout  string `yaml:"flush_timeout"`
	ChannelBuffer int    `yaml:"channel_buffer"`
}

// New 创建新的配置实例
func New() *Config {
	config := &Config{}

	// 首先尝试读取配置文件
	configPath := getEnv("CONFIG_PATH", "config.yaml")
	if err := loadConfigFromFile(config, configPath); err != nil {
		log.Printf("Warning: Failed to load config file %s: %v", configPath, err)
		log.Println("Falling back to environment variables and defaults")
	}

	// 使用环境变量覆盖配置文件设置（环境变量优先级更高）
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	} else if config.Port == "" {
		config.Port = "8080"
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.DBHost = dbHost
	} else if config.DBHost == "" {
		config.DBHost = "localhost"
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		config.DBPort = dbPort
	} else if config.DBPort == "" {
		config.DBPort = "3306"
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.DBUser = dbUser
	} else if config.DBUser == "" {
		config.DBUser = "root"
	}

	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.DBPassword = dbPassword
	} else if config.DBPassword == "" {
		config.DBPassword = ""
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.DBName = dbName
	} else if config.DBName == "" {
		config.DBName = "bench_server"
	}

	// 设置默认值
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 25
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 5
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.LogFormat == "" {
		config.LogFormat = "json"
	}
	if config.ReadTimeout == "" {
		config.ReadTimeout = "15s"
	}
	if config.WriteTimeout == "" {
		config.WriteTimeout = "15s"
	}
	if config.IdleTimeout == "" {
		config.IdleTimeout = "60s"
	}
	if config.BatchSize == 0 {
		config.BatchSize = 100
	}
	if config.FlushTimeout == "" {
		config.FlushTimeout = "5s"
	}
	if config.ChannelBuffer == 0 {
		config.ChannelBuffer = 10000
	}

	return config
}

// loadConfigFromFile 从文件加载配置
func loadConfigFromFile(config *Config, configPath string) error {
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析YAML配置文件
	var configFile ConfigFile
	if err := yaml.Unmarshal(data, &configFile); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 将解析的配置映射到Config结构体
	config.Port = configFile.Server.Port
	config.DBHost = configFile.Database.Host
	config.DBPort = configFile.Database.Port
	config.DBUser = configFile.Database.User
	config.DBPassword = configFile.Database.Password
	config.DBName = configFile.Database.Name
	config.MaxOpenConns = configFile.Database.MaxOpenConns
	config.MaxIdleConns = configFile.Database.MaxIdleConns
	config.LogLevel = configFile.Logging.Level
	config.LogFormat = configFile.Logging.Format
	config.ReadTimeout = configFile.App.ReadTimeout
	config.WriteTimeout = configFile.App.WriteTimeout
	config.IdleTimeout = configFile.App.IdleTimeout
	config.BatchSize = configFile.App.BatchSize
	config.FlushTimeout = configFile.App.FlushTimeout
	config.ChannelBuffer = configFile.App.ChannelBuffer

	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
