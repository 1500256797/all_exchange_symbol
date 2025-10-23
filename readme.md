# Exchange Symbol Tracker

一个用Go语言编写的加密货币交易所符号追踪系统，支持从多个主流交易所获取现货和期货交易对，检测新增符号并通过Telegram推送通知。

## 功能特性

- **多交易所支持**: 币安(Binance)、OKX、Gate.io、Bitget、Bybit
- **现货和期货**: 同时支持现货和期货交易对
- **自动检测**: 检测数据库中不存在的新符号
- **Telegram通知**: 自动推送新发现的符号到Telegram
- **数据库存储**: 使用MySQL存储符号信息
- **并发处理**: 高效的并发获取和处理

## 数据库结构

```go
type Symbol struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Exchange     string    `gorm:"not null;index" json:"exchange"`
    Type         string    `gorm:"not null;index" json:"type"` // "spot" or "futures"
    Symbol       string    `gorm:"not null;index" json:"symbol"`
    Combination  string    `gorm:"not null;unique" json:"combination"` // exchange-type-symbol
    CreatedAt    time.Time `json:"created_at"`
}
```

## 安装和使用

### 1. 克隆项目

```bash
git clone <repository-url>
cd all_exchange_symbol
```

### 2. 设置MySQL数据库

首先创建数据库：

```bash
# 登录MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE exchange_symbols CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 或者使用提供的SQL文件
mysql -u root -p < init.sql
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 配置环境变量

复制示例配置文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入你的配置：

```
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_mysql_password
MYSQL_DATABASE=exchange_symbols
LOG_LEVEL=info
```

### 5. 运行程序

```bash
# 从所有交易所获取符号
go run main.go

# 从指定交易所获取符号
go run main.go -exchange binance

# 查看数据库统计
go run main.go -stats

# 显示帮助信息
go run main.go -help
```

## 支持的交易所

| 交易所 | 现货API | 期货API |
|--------|---------|---------|
| Binance | ✅ | ✅ |
| OKX | ✅ | ✅ |
| Gate.io | ✅ | ✅ |
| Bitget | ✅ | ✅ |
| Bybit | ✅ | ✅ |

## 命令行选项

- `-exchange string`: 指定交易所 (binance, okx, gate, bitget, bybit)
- `-stats`: 显示数据库统计信息
- `-help`: 显示帮助信息

## Telegram Bot 设置

1. 在Telegram中找到 @BotFather
2. 发送 `/newbot` 创建新机器人
3. 按照指示设置机器人名称和用户名
4. 获取Bot Token并设置到环境变量
5. 获取Chat ID：
   - 将机器人添加到群组或私聊
   - 发送一条消息
   - 访问 `https://api.telegram.org/bot<YourBOTToken>/getUpdates`
   - 从响应中找到chat id

## 项目结构

```
all_exchange_symbol/
├── config/          # 配置管理
├── database/        # 数据库连接和初始化
├── exchanges/       # 各交易所API实现
├── models/          # 数据模型
├── processor/       # 数据处理逻辑
├── reader/          # 数据读取模块
├── writer/          # 数据写入和Telegram推送
├── main.go          # 主程序入口
├── go.mod           # Go模块文件
├── .env.example     # 环境变量示例
└── README.md        # 项目文档
```

## 工作流程

1. **Read**: 从支持的交易所并发获取现货和期货交易对
2. **Process**: 检查数据库中是否已存在该符号(基于组合键: exchange-type-symbol)
3. **Write**: 将新符号写入数据库并推送到Telegram

## 示例输出

```
2024/01/01 12:00:00 Starting exchange symbol synchronization...
2024/01/01 12:00:00 Fetching symbols from all exchanges
2024/01/01 12:00:00 MySQL database initialized successfully
2024/01/01 12:00:01 Successfully fetched 2000 spot symbols from binance
2024/01/01 12:00:01 Successfully fetched 500 futures symbols from binance
2024/01/01 12:00:02 Successfully fetched 1800 spot symbols from okx
...
2024/01/01 12:00:05 Fetched 10000 symbols in 5s
2024/01/01 12:00:05 New symbol found: binance-spot-NEWTOKEN
2024/01/01 12:00:05 Processed symbols in 100ms
2024/01/01 12:00:05 Successfully wrote 1 symbols to database
2024/01/01 12:00:05 Successfully sent message to Telegram with 1 new symbols
2024/01/01 12:00:05 Synchronization completed in 5.1s. Found 1 new symbols out of 10000 total.
```

## 注意事项

- 需要预先创建MySQL数据库，程序会自动创建表结构
- 确保MySQL用户有CREATE、SELECT、INSERT、UPDATE、DELETE权限
- 如果没有配置Telegram信息，程序仍会正常运行，只是不会发送通知
- 建议设置定时任务(如cron)定期运行程序检测新符号
- 各交易所的API可能有频率限制，程序已实现并发控制
- 建议为应用创建专用MySQL用户，不要使用root用户

## 定时任务设置

使用cron设置定时任务，例如每小时检查一次：

```bash
# 编辑crontab
crontab -e

# 添加以下行（每小时执行一次）
0 * * * * cd /path/to/all_exchange_symbol && go run main.go >> /var/log/exchange_symbols.log 2>&1
```