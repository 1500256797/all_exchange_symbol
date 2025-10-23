package exchanges

import (
	"all_exchange_symbol/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Binance struct {
	Name string
}

type BinanceSpotSymbol struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`
}

type BinanceFuturesSymbol struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`
}

func NewBinance() *Binance {
	return &Binance{Name: "binance"}
}

func (b *Binance) GetName() string {
	return b.Name
}

func (b *Binance) FetchSpotSymbols() ([]models.Symbol, error) {
	log.Printf("开始获取币安现货交易对数据...")

	resp, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
	if err != nil {
		log.Printf("币安现货API请求失败: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("币安现货API请求成功，状态码: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取币安现货API响应失败: %v", err)
		return nil, err
	}

	var result struct {
		Symbols []BinanceSpotSymbol `json:"symbols"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("解析币安现货API响应失败: %v", err)
		return nil, err
	}

	log.Printf("从币安API获取到 %d 个现货交易对", len(result.Symbols))

	var symbols []models.Symbol

	for _, s := range result.Symbols {
		symbols = append(symbols, models.Symbol{
			Exchange:  b.Name,
			Type:      "spot",
			Symbol:    s.Symbol,
			CreatedAt: time.Now(),
		})
	}

	log.Printf("币安现货交易对处理完成 - 共 %d 个", len(symbols))
	return symbols, nil
}

func (b *Binance) FetchFuturesSymbols() ([]models.Symbol, error) {
	log.Printf("开始获取币安合约交易对数据...")

	resp, err := http.Get("https://fapi.binance.com/fapi/v1/exchangeInfo")
	if err != nil {
		log.Printf("币安合约API请求失败: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("币安合约API请求成功，状态码: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取币安合约API响应失败: %v", err)
		return nil, err
	}

	var result struct {
		Symbols []BinanceFuturesSymbol `json:"symbols"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("解析币安合约API响应失败: %v", err)
		return nil, err
	}

	log.Printf("从币安API获取到 %d 个合约交易对", len(result.Symbols))

	var symbols []models.Symbol

	for _, s := range result.Symbols {
		symbols = append(symbols, models.Symbol{
			Exchange:  b.Name,
			Type:      "futures",
			Symbol:    s.Symbol,
			CreatedAt: time.Now(),
		})
	}

	log.Printf("币安合约交易对处理完成 - 共 %d 个", len(symbols))
	return symbols, nil
}
