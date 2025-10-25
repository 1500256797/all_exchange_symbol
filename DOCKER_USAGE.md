# Docker 使用说明

## 快速开始

### 1. 配置环境变量（可选）

如果需要使用 Telegram 通知功能，可以创建 `.env` 文件：

```bash
# 复制示例配置
cat > .env << 'EOF'
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
TELEGRAM_CHAT_ID=your_telegram_chat_id_here
EOF
```

### 2. 启动所有服务

```bash
# 构建并启动服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 只查看应用日志
docker-compose logs -f app

# 只查看 MySQL 日志
docker-compose logs -f mysql
```

### 3. 停止服务

```bash
# 停止服务
docker-compose down

# 停止服务并删除数据卷（注意：会删除数据库数据）
docker-compose down -v
```

## 服务说明

### MySQL 服务

-   **容器名**: `exchange_symbols_mysql`
-   **端口**: 3306
-   **用户名**: root
-   **密码**: your_mysql_password
-   **数据库**: exchange_symbols
-   **数据持久化**: 使用 Docker volume `mysql_data`

### Go 应用服务

-   **容器名**: `exchange_symbols_app`
-   **运行模式**: daemon 模式（每 5 秒检查一次）
-   **自动重启**: 是
-   **依赖**: 等待 MySQL 健康检查通过后启动

## 常用命令

```bash
# 重新构建应用
docker-compose up -d --build app

# 查看应用实时日志
docker-compose logs -f app

# 进入 MySQL 容器
docker-compose exec mysql mysql -uroot -pyour_mysql_password exchange_symbols

# 进入应用容器
docker-compose exec app sh

# 重启应用
docker-compose restart app

# 重启所有服务
docker-compose restart

# 查看资源使用情况
docker stats exchange_symbols_app exchange_symbols_mysql
```

## 数据库操作

### 连接数据库

```bash
# 从宿主机连接
mysql -h 127.0.0.1 -P 3306 -uroot -pyour_mysql_password exchange_symbols

# 或使用 docker-compose
docker-compose exec mysql mysql -uroot -pyour_mysql_password exchange_symbols
```

### 备份数据库

```bash
# 导出数据库
docker-compose exec mysql mysqldump -uroot -pyour_mysql_password exchange_symbols > backup_$(date +%Y%m%d_%H%M%S).sql

# 导入数据库
docker-compose exec -T mysql mysql -uroot -pyour_mysql_password exchange_symbols < backup.sql
```

### 查看数据统计

```bash
# 进入容器执行统计命令
docker-compose exec app ./main -stats
```

## 自定义配置

### 修改 MySQL 密码

编辑 `docker-compose.yml` 文件，修改以下位置的密码：

1. MySQL 服务的 `MYSQL_ROOT_PASSWORD`
2. MySQL 服务的 healthcheck 中的密码
3. App 服务的 `MYSQL_PASSWORD`

### 修改检查间隔

编辑 `main.go` 第 105 行：

```go
ticker := time.NewTicker(5 * time.Second)  // 改为你想要的间隔
```

然后重新构建：

```bash
docker-compose up -d --build app
```

### 只监控特定交易所

修改 `docker-compose.yml` 中的 CMD：

```yaml
# 在 Dockerfile 中修改最后一行
CMD ["./main", "-daemon", "-exchange", "binance"]
```

## 故障排查

### 应用无法连接数据库

```bash
# 检查 MySQL 是否健康
docker-compose ps

# 查看 MySQL 日志
docker-compose logs mysql

# 测试数据库连接
docker-compose exec mysql mysql -uroot -pyour_mysql_password -e "SELECT 1"
```

### 查看应用错误日志

```bash
# 查看最近的日志
docker-compose logs --tail=100 app

# 实时跟踪日志
docker-compose logs -f app
```

### 重置所有数据

```bash
# 停止并删除所有容器、网络、数据卷
docker-compose down -v

# 重新启动
docker-compose up -d --build
```

## 生产环境建议

1. **修改默认密码**: 将 `your_mysql_password` 改为强密码
2. **环境变量管理**: 使用 `.env` 文件管理敏感信息
3. **日志管理**: 配置日志轮转，避免日志文件过大
4. **监控**: 添加监控和告警机制
5. **备份**: 定期备份数据库数据
6. **资源限制**: 在 docker-compose.yml 中添加资源限制

```yaml
app:
    # ...其他配置
    deploy:
        resources:
            limits:
                cpus: '1'
                memory: 512M
            reservations:
                cpus: '0.5'
                memory: 256M
```
