package main

import (
	"all_exchange_symbol/config"
	"all_exchange_symbol/database"
	"all_exchange_symbol/models"
	"all_exchange_symbol/processor"
	"all_exchange_symbol/reader"
	"all_exchange_symbol/writer"
	"flag"
	"log"
	"time"
)

func main() {
	var (
		exchangeFlag = flag.String("exchange", "", "Fetch symbols from specific exchange (binance, okx, gate, bitget, bybit)")
		helpFlag     = flag.Bool("help", false, "Show help information")
		statsFlag    = flag.Bool("stats", false, "Show database statistics")
		verifyFlag   = flag.Bool("verify", false, "Compare API data with database data for detailed verification")
		daemonFlag   = flag.Bool("daemon", false, "Run in daemon mode with 5-second periodic checks")
	)
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	cfg := config.Load()

	database.Initialize()
	defer database.Close()

	if *statsFlag {
		showStats()
		return
	}

	if *verifyFlag {
		showDataVerification(*exchangeFlag)
		return
	}

	if *daemonFlag {
		runDaemon(*exchangeFlag, cfg)
		return
	}

	log.Println("Starting exchange symbol synchronization...")
	start := time.Now()

	r := reader.NewReader()
	p := processor.NewProcessor()
	w := writer.NewWriter(cfg.TelegramBotToken, cfg.TelegramChatID)

	var fetchedSymbols []models.Symbol
	var err error

	if *exchangeFlag != "" {
		log.Printf("Fetching symbols from %s only", *exchangeFlag)
		fetchedSymbols, err = r.FetchSymbolsByExchange(*exchangeFlag)
	} else {
		log.Println("Fetching symbols from all exchanges")
		fetchedSymbols, err = r.FetchAllSymbols()
	}

	if err != nil {
		log.Fatalf("Error fetching symbols: %v", err)
	}

	log.Printf("Fetched %d symbols in %v", len(fetchedSymbols), time.Since(start))

	processStart := time.Now()
	newSymbols, err := p.ProcessSymbols(fetchedSymbols)
	if err != nil {
		log.Fatalf("Error processing symbols: %v", err)
	}

	log.Printf("Processed symbols in %v", time.Since(processStart))

	writeStart := time.Now()
	if err := w.ProcessAndWrite(newSymbols); err != nil {
		log.Fatalf("Error writing symbols: %v", err)
	}

	log.Printf("Wrote symbols in %v", time.Since(writeStart))

	if err := w.SendSummaryToTelegram(len(fetchedSymbols), len(newSymbols)); err != nil {
		log.Printf("Error sending summary: %v", err)
	}

	log.Printf("Synchronization completed in %v. Found %d new symbols out of %d total.",
		time.Since(start), len(newSymbols), len(fetchedSymbols))
}

func runDaemon(exchange string, cfg *config.Config) {
	log.Println("Starting daemon mode with 5-second intervals...")
	if exchange != "" {
		log.Printf("Monitoring exchange: %s", exchange)
	} else {
		log.Println("Monitoring all exchanges")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			performSynchronization(exchange, cfg)
		}
	}
}

func performSynchronization(exchange string, cfg *config.Config) {
	start := time.Now()
	log.Printf("[%s] Starting synchronization check...", start.Format("15:04:05"))

	r := reader.NewReader()
	p := processor.NewProcessor()
	w := writer.NewWriter(cfg.TelegramBotToken, cfg.TelegramChatID)

	var fetchedSymbols []models.Symbol
	var err error

	if exchange != "" {
		fetchedSymbols, err = r.FetchSymbolsByExchange(exchange)
	} else {
		fetchedSymbols, err = r.FetchAllSymbols()
	}

	if err != nil {
		log.Printf("Error fetching symbols: %v", err)
		return
	}

	newSymbols, err := p.ProcessSymbols(fetchedSymbols)
	if err != nil {
		log.Printf("Error processing symbols: %v", err)
		return
	}

	if err := w.ProcessAndWrite(newSymbols); err != nil {
		log.Printf("Error writing symbols: %v", err)
		return
	}

	if len(newSymbols) > 0 {
		log.Printf("[%s] Found %d new symbols out of %d total (took %v)",
			start.Format("15:04:05"), len(newSymbols), len(fetchedSymbols), time.Since(start))

		if err := w.SendSummaryToTelegram(len(fetchedSymbols), len(newSymbols)); err != nil {
			log.Printf("Error sending summary: %v", err)
		}
	} else {
		log.Printf("[%s] No new symbols found (%d symbols checked, took %v)",
			start.Format("15:04:05"), len(fetchedSymbols), time.Since(start))
	}
}

func showHelp() {
	log.Println(`
Exchange Symbol Synchronizer

Usage:
  go run main.go [options]

Options:
  -exchange string    Fetch symbols from specific exchange (binance, okx, gate, bitget, bybit)
  -stats              Show database statistics
  -verify             Compare API data with database data for detailed verification
  -daemon             Run in daemon mode with 5-second periodic checks
  -help               Show this help message

Examples:
  go run main.go                        # Fetch from all exchanges
  go run main.go -exchange binance      # Fetch from Binance only
  go run main.go -stats                 # Show database statistics
  go run main.go -verify                # Verify API vs database for all exchanges
  go run main.go -verify -exchange binance # Verify API vs database for Binance only
  go run main.go -daemon                # Run daemon mode checking every 5 seconds
  go run main.go -daemon -exchange binance # Run daemon mode for Binance only

Environment Variables:
  TELEGRAM_BOT_TOKEN    Your Telegram bot token
  TELEGRAM_CHAT_ID      Your Telegram chat ID
  DATABASE_PATH         Database file path (default: symbols.db)
  LOG_LEVEL             Log level (default: info)
`)
}

func showStats() {
	p := processor.NewProcessor()

	total, err := p.GetSymbolCount()
	if err != nil {
		log.Fatalf("Error getting total count: %v", err)
	}

	log.Printf("Total symbols in database: %d", total)

	exchanges := []string{"binance", "okx", "gate", "bitget", "bybit"}
	for _, exchange := range exchanges {
		count, err := p.GetSymbolCountByExchange(exchange)
		if err != nil {
			log.Printf("Error getting count for %s: %v", exchange, err)
			continue
		}
		log.Printf("  %s: %d symbols", exchange, count)
	}

	spotCount, err := p.GetExistingSymbolsByType("spot")
	if err != nil {
		log.Printf("Error getting spot count: %v", err)
	} else {
		log.Printf("Spot symbols: %d", len(spotCount))
	}

	futuresCount, err := p.GetExistingSymbolsByType("futures")
	if err != nil {
		log.Printf("Error getting futures count: %v", err)
	} else {
		log.Printf("Futures symbols: %d", len(futuresCount))
	}
}

func showDataVerification(exchange string) {
	log.Println("=== 开始API与数据库数据验证 ===")

	r := reader.NewReader()
	p := processor.NewProcessor()

	start := time.Now()

	var exchanges []string
	if exchange != "" {
		exchanges = []string{exchange}
		log.Printf("验证交易所: %s", exchange)
	} else {
		exchanges = []string{"binance", "okx", "gate", "bitget", "bybit"}
		log.Println("验证所有交易所")
	}

	for _, ex := range exchanges {
		log.Printf("\n=== 开始验证 %s 交易所数据 ===", ex)

		fetchedSymbols, err := r.FetchSymbolsByExchange(ex)
		if err != nil {
			log.Printf("获取%s交易所数据失败: %v", ex, err)
			continue
		}

		if len(fetchedSymbols) == 0 {
			log.Printf("未获取到%s交易所数据", ex)
			continue
		}

		var spotSymbols, futuresSymbols []models.Symbol
		for _, symbol := range fetchedSymbols {
			if symbol.Type == "spot" {
				spotSymbols = append(spotSymbols, symbol)
			} else if symbol.Type == "futures" {
				futuresSymbols = append(futuresSymbols, symbol)
			}
		}

		if len(spotSymbols) > 0 {
			log.Printf("\n--- 验证 %s 现货数据 ---", ex)
			_, err := p.CompareAPIWithDatabase(spotSymbols, ex, "spot")
			if err != nil {
				log.Printf("现货数据对比失败: %v", err)
			}
		}

		if len(futuresSymbols) > 0 {
			log.Printf("\n--- 验证 %s 合约数据 ---", ex)
			_, err := p.CompareAPIWithDatabase(futuresSymbols, ex, "futures")
			if err != nil {
				log.Printf("合约数据对比失败: %v", err)
			}
		}

		log.Printf("=== %s 交易所验证完成 ===", ex)
	}

	log.Printf("\n=== 全部验证完成，耗时: %v ===", time.Since(start))
}
