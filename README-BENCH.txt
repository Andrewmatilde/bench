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

数据配置：
  key_range           int      设备ID范围 (默认: 1000)
  report_interval     int      实时报告间隔（秒）(默认: 5)

MySQL配置：
  mysql_dsn           string   MySQL数据源名称 (默认: "")

上报配置：
  report_url          string   统计数据上报URL (默认: "")
  report_key          string   上报认证密钥，请填写自己组的具体组名。必填参数，用于区分各小组上报身份 (默认: ""，可选参数为 team1/team2/team3/team4/team5)

示例配置文件 (config.json)：
{
  "server_url": "http://localhost:8080",
  "duration_seconds": 60,
  "mode": "qps",
  "qps": 100,
  "key_range": 1000,
  "report_interval": 5,
  "mysql_dsn": "",
  "report_url": "http://monitoring-server/api/stats",
  "report_key": "your-team-key"  // 将同时设置 X-Team-ID 和 X-Team-Name header
}

使用方法：
  ./client -config config.json
  ./client -help-config