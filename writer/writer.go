package writer

import (
	"all_exchange_symbol/database"
	"all_exchange_symbol/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Writer struct {
	telegramBotToken string
	telegramChatID   string
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func NewWriter(botToken, chatID string) *Writer {
	return &Writer{
		telegramBotToken: botToken,
		telegramChatID:   chatID,
	}
}

func (w *Writer) WriteSymbolsToDatabase(symbols []models.Symbol) error {
	if len(symbols) == 0 {
		log.Println("No new symbols to write to database")
		return nil
	}

	result := database.DB.Create(&symbols)
	if result.Error != nil {
		log.Printf("Error writing symbols to database: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully wrote %d symbols to database", len(symbols))
	return nil
}

func (w *Writer) SendToTelegram(symbols []models.Symbol) error {
	if len(symbols) == 0 {
		log.Println("No new symbols to send to Telegram")
		return nil
	}

	message := w.formatTelegramMessage(symbols)

	telegramMsg := TelegramMessage{
		ChatID:    w.telegramChatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		log.Printf("Error marshaling telegram message: %v", err)
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", w.telegramBotToken)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending telegram message: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API returned status code: %d", resp.StatusCode)
		return fmt.Errorf("telegram API error: status code %d", resp.StatusCode)
	}

	log.Printf("Successfully sent message to Telegram with %d new symbols", len(symbols))
	return nil
}

func (w *Writer) formatTelegramMessage(symbols []models.Symbol) string {
	if len(symbols) == 0 {
		return "No new symbols found."
	}

	message := fmt.Sprintf("🚀 *Found %d new trading symbols:*\n\n", len(symbols))

	exchangeGroups := make(map[string][]models.Symbol)
	for _, symbol := range symbols {
		exchangeGroups[symbol.Exchange] = append(exchangeGroups[symbol.Exchange], symbol)
	}

	for exchange, exchangeSymbols := range exchangeGroups {
		spotCount := 0
		futuresCount := 0

		for _, symbol := range exchangeSymbols {
			if symbol.Type == "spot" {
				spotCount++
			} else {
				futuresCount++
			}
		}

		message += fmt.Sprintf("📊 *%s*:\n", exchange)
		if spotCount > 0 {
			message += fmt.Sprintf("   • Spot: %d symbols\n", spotCount)
		}
		if futuresCount > 0 {
			message += fmt.Sprintf("   • Futures: %d symbols\n", futuresCount)
		}

		if len(exchangeSymbols) <= 10 {
			for _, symbol := range exchangeSymbols {
				message += fmt.Sprintf("   - `%s`\n", symbol.Symbol)
			}
		} else {
			for i, symbol := range exchangeSymbols[:5] {
				message += fmt.Sprintf("   - `%s`\n", symbol.Symbol)
				if i == 4 {
					message += fmt.Sprintf("   ... and %d more\n", len(exchangeSymbols)-5)
				}
			}
		}
		message += "\n"
	}

	return message
}

func (w *Writer) ProcessAndWrite(symbols []models.Symbol) error {
	if err := w.WriteSymbolsToDatabase(symbols); err != nil {
		return fmt.Errorf("failed to write to database: %v", err)
	}

	if w.telegramBotToken != "" && w.telegramChatID != "" {
		if err := w.SendToTelegram(symbols); err != nil {
			log.Printf("Failed to send to Telegram (continuing anyway): %v", err)
		}
	} else {
		log.Println("Telegram credentials not provided, skipping notification")
	}

	return nil
}

func (w *Writer) SendSummaryToTelegram(totalSymbols int, newSymbols int) error {
	if w.telegramBotToken == "" || w.telegramChatID == "" {
		log.Println("Telegram credentials not provided, skipping summary")
		return nil
	}

	message := "📈 *Symbol Sync Summary*\n\n"
	message += fmt.Sprintf("🔍 Total symbols checked: %d\n", totalSymbols)
	message += fmt.Sprintf("✨ New symbols found: %d\n", newSymbols)

	if newSymbols == 0 {
		message += "\n✅ No new symbols detected. All markets are up to date!"
	}

	telegramMsg := TelegramMessage{
		ChatID:    w.telegramChatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", w.telegramBotToken)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: status code %d", resp.StatusCode)
	}

	log.Println("Summary sent to Telegram successfully")
	return nil
}
