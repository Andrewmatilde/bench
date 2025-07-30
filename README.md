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

## 📝 许可证

MIT License 