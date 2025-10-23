-- 创建数据库 (如果不存在)
CREATE DATABASE IF NOT EXISTS exchange_symbols 
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE exchange_symbols;

-- 创建符号表 (GORM会自动创建表，这里只是展示表结构)
-- CREATE TABLE IF NOT EXISTS symbols (
--   id BIGINT AUTO_INCREMENT PRIMARY KEY,
--   exchange VARCHAR(50) NOT NULL,
--   type VARCHAR(20) NOT NULL,
--   symbol VARCHAR(100) NOT NULL,
--   combination VARCHAR(200) NOT NULL UNIQUE,
--   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--   INDEX idx_exchange (exchange),
--   INDEX idx_type (type),
--   INDEX idx_symbol (symbol),
--   INDEX idx_combination (combination)
-- );