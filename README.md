# Bench Server

一个高性能的Go传感器数据处理服务器，提供传感器数据上报、读写操作、批量处理和统计查询等功能。支持配置文件和环境变量两种配置方式。

## 🏗️ 项目结构

```
bench-server/
├── main.go                 # 主程序入口
├── config.yaml            # 配置文件
├── pkg/                   # 包目录
│   ├── config/           # 配置管理
│   │   └── config.go
│   ├── database/         # 数据库操作
│   │   └── database.go
│   ├── handlers/         # HTTP处理器
│   │   └── handlers.go
│   ├── utils/           # 工具函数
│   │   └── utils.go
│   └── writer/          # 写入器（批量、压缩、优先级）
│       └── writer.go
├── init.sql             # 数据库初始化脚本
├── go.mod              # Go模块文件
├── go.sum              # 依赖校验文件
└── docs/               # 文档目录
    ├── openapi.yaml    # OpenAPI规范
    └── ...
```

## 🚀 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+

### 安装和运行

1. **克隆仓库**
```bash
git clone <repository-url>
cd bench-server
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置数据库**
```bash
# 确保MySQL服务运行
mysql -u root -p < init.sql
```

4. **配置应用**

方式1: 使用配置文件 (推荐)
```yaml
# config.yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: "3306"
  user: "root"
  password: "your_password"
  name: "bench_server"
  max_open_conns: 25
  max_idle_conns: 5

logging:
  level: "info"
  format: "json"

app:
  read_timeout: "15s"
  write_timeout: "15s"
  idle_timeout: "60s"
```

方式2: 使用环境变量
```bash
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_USER="root"
export DB_PASSWORD="your_password"
export DB_NAME="bench_server"
export PORT="8080"
```

5. **编译和运行**
```bash
go build -o bench-server
./bench-server
```

## 📦 包说明

### pkg/config
配置管理包，支持：
- YAML配置文件解析
- 环境变量覆盖
- 配置验证和默认值

### pkg/database
数据库操作包，提供：
- 数据库连接管理
- 表结构初始化
- 传感器数据CRUD操作
- 统计查询功能

### pkg/handlers
HTTP处理器包，包含：
- RESTful API端点
- 请求验证和错误处理
- JSON响应格式化

### pkg/utils
工具函数包，提供：
- 随机负载数据生成
- Base64编码/解码
- 数据验证工具

### pkg/writer
高级写入器包，支持：
- 批量写入优化
- 数据压缩
- 优先级队列
- 异步处理

## 🔌 API 端点

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查 |
| POST | `/api/sensor-data` | 传感器数据上报 |
| POST | `/api/sensor-rw` | 传感器读写操作 |
| POST | `/api/batch-sensor-rw` | 批量传感器读写 |
| GET | `/api/stats` | 统计信息查询 |
| POST | `/api/get-sensor-data` | 传感器数据查询 |

### 示例API调用

**传感器数据上报**
```bash
curl -X POST http://localhost:8080/api/sensor-data \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2024-01-15T10:30:00.000Z",
    "device_id": "device001",
    "metric_name": "temperature",
    "value": 25.8,
    "priority": 1,
    "data": "sensor payload data"
  }'
```

**统计信息查询**
```bash
curl http://localhost:8080/api/stats
```

## ⚙️ 配置选项

### 配置优先级
1. **环境变量** (最高优先级)
2. **配置文件** (中等优先级)
3. **默认值** (最低优先级)

### 环境变量列表
- `CONFIG_PATH`: 配置文件路径 (默认: config.yaml)
- `PORT`: 服务器端口
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称

## 🗄️ 数据库架构

### 主要表结构

**time_series_data** - 时序数据表
- 存储传感器时序数据
- 支持毫秒级时间戳
- 包含设备ID、指标名称、数值、优先级等字段

**device_status** - 设备状态表
- 存储设备当前状态
- 记录最后更新时间和告警计数

## 🔧 开发指南

### 项目架构

采用分层架构设计：
```
main.go -> pkg/handlers -> pkg/database
  ↓            ↓              ↓
config     HTTP处理        数据持久化
  ↓            ↓              ↓
pkg/config  中间件          SQL操作
            验证
```

### 添加新功能

1. **新增API端点**：在 `pkg/handlers` 中添加处理器
2. **数据库操作**：在 `pkg/database` 中添加相关方法
3. **工具函数**：在 `pkg/utils` 中添加通用工具
4. **配置选项**：在 `pkg/config` 中扩展配置结构

### 代码规范

- 遵循Go官方代码规范
- 使用有意义的包名和函数名
- 添加适当的注释和文档
- 错误处理要完整

## 🧪 测试

### 运行测试
```bash
go test ./...
```

### API测试
项目包含完整的API测试脚本，详见测试目录。

## 🚧 负载测试

项目包含完整的负载测试工具，支持：
- 并发请求测试
- 吞吐量基准测试
- 延迟分析
- 错误率统计

### ts-bench 基准测试工具

项目提供了专业的时序基准测试工具 `ts-bench`，用于全面测试bench-server的性能。

#### 主要特性

- **多种流量控制模式**: 支持QPS限制和并发数控制
- **操作类型比例配置**: 可自定义各种API操作的比例
- **实时性能监控**: 提供实时性能指标和统计报告
- **MySQL直连统计**: 支持直接从数据库获取统计信息
- **远程数据上报**: 支持将测试结果上报到监控系统

#### 使用方法

```bash
# 使用默认配置
./ts-bench

# 使用自定义配置文件
./ts-bench -config custom-config.json

# 查看配置说明
./ts-bench -help-config
```

#### 配置文件示例

```json
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
  "mysql_dsn": "root:password@tcp(localhost:3306)/bench_server",
  "report_url": "http://monitoring-server/api/stats",
  "report_key": "your-team-key"
}
```

#### 配置参数说明

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `server_url` | string | 目标服务器地址 | http://localhost:8080 |
| `duration_seconds` | int | 测试持续时间（秒） | 60 |
| `mode` | string | 流量控制模式: "qps" 或 "concurrency" | qps |
| `qps` | int | 目标QPS（QPS模式） | 100 |
| `concurrency` | int | 并发数（并发模式） | 10 |
| `sensor_data_ratio` | float64 | 传感器数据上报比例 | 0.4 |
| `sensor_rw_ratio` | float64 | 传感器读写操作比例 | 0.3 |
| `batch_rw_ratio` | float64 | 批量操作比例 | 0.2 |
| `query_ratio` | float64 | 查询操作比例 | 0.1 |
| `key_range` | int | 设备ID范围 | 1000 |
| `report_interval` | int | 实时报告间隔（秒） | 5 |
| `mysql_dsn` | string | MySQL数据源名称 | "" |
| `report_url` | string | 统计数据上报URL | "" |
| `report_key` | string | 上报认证密钥 | "" |

#### 测试场景

**基础性能测试**
```bash
# 100 QPS 持续60秒
./ts-bench -config basic-test.json
```

**高并发测试**
```json
{
  "mode": "concurrency",
  "concurrency": 50,
  "duration_seconds": 300
}
```

**写入密集型测试**
```json
{
  "sensor_data_ratio": 0.6,
  "sensor_rw_ratio": 0.3,
  "batch_rw_ratio": 0.1,
  "query_ratio": 0.0
}
```

**查询密集型测试**
```json
{
  "sensor_data_ratio": 0.1,
  "sensor_rw_ratio": 0.1,
  "batch_rw_ratio": 0.1,
  "query_ratio": 0.7
}
```

#### 性能指标说明

测试工具会实时输出以下关键指标：

- **QPS**: 每秒请求数 (Queries Per Second)
- **平均延迟**: 请求平均响应时间
- **P95/P99延迟**: 95%/99%请求的响应时间
- **错误率**: 失败请求占总请求的比例
- **吞吐量**: 数据处理速度
- **数据库指标**: 连接数、查询性能等

#### 测试最佳实践

1. **预热测试**: 先进行短时间低负载测试
2. **渐进式压测**: 逐步增加负载找到性能瓶颈
3. **长时间稳定性测试**: 验证系统在持续负载下的稳定性
4. **监控系统资源**: 关注CPU、内存、磁盘IO等指标
5. **数据库优化**: 根据测试结果调整数据库配置

#### 示例测试计划

```bash
# 1. 基础功能验证 (低负载)
./ts-bench -config config-basic.json

# 2. 性能基准测试 (中等负载)  
./ts-bench -config config-benchmark.json

# 3. 压力测试 (高负载)
./ts-bench -config config-stress.json

# 4. 稳定性测试 (长时间)
./ts-bench -config config-stability.json
```

## 📝 许可证

MIT License 