时序基准测试工具: 
ts-bench
  -config string
        配置文件路径 (default "config.json")
  -help-config
        显示配置结构说明

=== 配置结构说明 ===

配置文件使用 JSON 格式，包含以下字段：

服务器配置：
  server_url          string   服务器地址 (默认: http://localhost:8080)
  duration_seconds    int      测试持续时间（秒）(默认: 60)

流量控制配置：
  mode                string   流量控制模式: "qps" 或 "concurrency" (默认: qps)
  qps                 int      目标QPS（mode=qps时使用）(默认: 100)
  concurrency         int      并发数（mode=concurrency时使用）(默认: 10)

操作比例配置（总和应≤1.0）：
  sensor_data_ratio   float64  传感器数据上报比例 (默认: 0.4)
  sensor_rw_ratio     float64  传感器读写操作比例 (默认: 0.3)
  batch_rw_ratio      float64  批量操作比例 (默认: 0.2)
  query_ratio         float64  查询操作比例 (默认: 0.1)

数据配置：
  key_range           int      设备ID范围 (默认: 1000)
  report_interval     int      实时报告间隔（秒）(默认: 5)

MySQL配置：
  mysql_dsn           string   MySQL数据源名称 (默认: "")

上报配置：
  report_url          string   统计数据上报URL (默认: "")
  report_key          string   上报认证密钥，用于设置 X-Team-ID 和 X-Team-Name header (默认: "")

示例配置文件 (config.json)：
{
  "server_url": "http://localhost:8080",
  "duration_seconds": 60,
  "mode": "qps",
  "qps": 100,
  "sensor_data_ratio": 0.4,
  "sensor_rw_ratio": 0.3,
  "batch_rw_ratio": 0.2,
  "query_ratio": 0.1,
  "key_range": 1000,
  "report_interval": 5,
  "mysql_dsn": "",
  "report_url": "http://monitoring-server/api/stats",
  "report_key": "your-team-key"  // 将同时设置 X-Team-ID 和 X-Team-Name header
}

使用方法：
  ./client -config config.json
  ./client -help-config