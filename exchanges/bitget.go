package exchanges

import (
	"all_exchange_symbol/models"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Bitget struct {
	Name string
}

type BitgetSymbol struct {
	Symbol            string `json:"symbol"`
	BaseCoin          string `json:"baseCoin"`
	QuoteCoin         string `json:"quoteCoin"`
	MinTradeAmount    string `json:"minTradeAmount"`
	MaxTradeAmount    string `json:"maxTradeAmount"`
	TakerFeeRate      string `json:"takerFeeRate"`
	MakerFeeRate      string `json:"makerFeeRate"`
	PricePrecision    string `json:"pricePrecision"`
	QuantityPrecision string `json:"quantityPrecision"`
	Status            string `json:"status"`
}

type BitgetFuturesSymbol struct {
	Symbol              string `json:"symbol"`
	BaseCoin            string `json:"baseCoin"`
	QuoteCoin           string `json:"quoteCoin"`
	BuyLimitPriceRatio  string `json:"buyLimitPriceRatio"`
	SellLimitPriceRatio string `json:"sellLimitPriceRatio"`
	FeeRateUpRatio      string `json:"feeRateUpRatio"`
	MakerFeeRate        string `json:"makerFeeRate"`
	TakerFeeRate        string `json:"takerFeeRate"`
	Status              string `json:"status"`
}

func NewBitget() *Bitget {
	return &Bitget{Name: "bitget"}
}

func (b *Bitget) GetName() string {
	return b.Name
}

func (b *Bitget) FetchSpotSymbols() ([]models.Symbol, error) {
	resp, err := http.Get("https://api.bitget.com/api/spot/v1/public/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code string         `json:"code"`
		Msg  string         `json:"msg"`
		Data []BitgetSymbol `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	for _, s := range result.Data {
		symbols = append(symbols, models.Symbol{
			Exchange:  b.Name,
			Type:      "spot",
			Symbol:    s.Symbol,
			CreatedAt: time.Now(),
		})
	}

	return symbols, nil
}

func (b *Bitget) FetchFuturesSymbols() ([]models.Symbol, error) {
	resp, err := http.Get("https://api.bitget.com/api/mix/v1/market/contracts?productType=umcbl")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code string                `json:"code"`
		Msg  string                `json:"msg"`
		Data []BitgetFuturesSymbol `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	for _, s := range result.Data {
		symbols = append(symbols, models.Symbol{
			Exchange:  b.Name,
			Type:      "futures",
			Symbol:    s.Symbol,
			CreatedAt: time.Now(),
		})
	}

	return symbols, nil
}
