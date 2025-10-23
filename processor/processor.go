package processor

import (
	"all_exchange_symbol/database"
	"all_exchange_symbol/models"
	"log"
	"sort"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) ProcessSymbols(fetchedSymbols []models.Symbol) ([]models.Symbol, error) {
	log.Printf("=== 开始处理交易对数据 ===")
	log.Printf("从API获取的交易对总数: %d", len(fetchedSymbols))

	var newSymbols []models.Symbol
	var existingCount int

	exchangeCounts := make(map[string]map[string]int)

	// 批量获取所有现有的交易对组合以提高性能
	log.Printf("正在批量获取现有交易对数据...")
	existingSymbols, err := p.GetAllExistingSymbols()
	if err != nil {
		return nil, err
	}

	// 创建一个map用于快速查找现有的交易对
	existingCombinations := make(map[string]bool)
	for _, symbol := range existingSymbols {
		combination := symbol.Exchange + "-" + symbol.Type + "-" + symbol.Symbol
		existingCombinations[combination] = true
	}
	log.Printf("数据库中现有交易对总数: %d", len(existingSymbols))

	for _, symbol := range fetchedSymbols {
		if _, ok := exchangeCounts[symbol.Exchange]; !ok {
			exchangeCounts[symbol.Exchange] = make(map[string]int)
		}
		exchangeCounts[symbol.Exchange][symbol.Type]++

		combination := symbol.Exchange + "-" + symbol.Type + "-" + symbol.Symbol
		exists := existingCombinations[combination]

		if !exists {
			newSymbols = append(newSymbols, symbol)
			log.Printf("发现新交易对: %s-%s-%s", symbol.Exchange, symbol.Type, symbol.Symbol)
		} else {
			existingCount++
		}
	}

	log.Printf("\n=== 数据处理统计 ===")
	for exchange, types := range exchangeCounts {
		log.Printf("交易所 %s:", exchange)
		for symbolType, count := range types {
			log.Printf("  %s: %d 个", symbolType, count)
		}
	}

	log.Printf("\n数据库中已存在: %d 个", existingCount)
	log.Printf("新发现的交易对: %d 个", len(newSymbols))
	log.Printf("处理完成，总处理: %d 个", len(fetchedSymbols))

	return newSymbols, nil
}

func (p *Processor) CheckSymbolExists(symbol models.Symbol) (bool, error) {
	var existingSymbol models.Symbol

	result := database.DB.Where("combination = ?", symbol.Exchange+"-"+symbol.Type+"-"+symbol.Symbol).First(&existingSymbol)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func (p *Processor) GetExistingSymbolsByExchange(exchange string) ([]models.Symbol, error) {
	var symbols []models.Symbol

	result := database.DB.Where("exchange = ?", exchange).Find(&symbols)
	if result.Error != nil {
		return nil, result.Error
	}

	return symbols, nil
}

func (p *Processor) GetExistingSymbolsByType(symbolType string) ([]models.Symbol, error) {
	var symbols []models.Symbol

	result := database.DB.Where("type = ?", symbolType).Find(&symbols)
	if result.Error != nil {
		return nil, result.Error
	}

	return symbols, nil
}

func (p *Processor) GetAllExistingSymbols() ([]models.Symbol, error) {
	var symbols []models.Symbol

	result := database.DB.Find(&symbols)
	if result.Error != nil {
		return nil, result.Error
	}

	return symbols, nil
}

func (p *Processor) GetSymbolCount() (int64, error) {
	var count int64

	result := database.DB.Model(&models.Symbol{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (p *Processor) GetSymbolCountByExchange(exchange string) (int64, error) {
	var count int64

	result := database.DB.Model(&models.Symbol{}).Where("exchange = ?", exchange).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (p *Processor) getSymbolsByExchangeAndType(exchange, symbolType string) ([]models.Symbol, error) {
	var symbols []models.Symbol

	result := database.DB.Where("exchange = ? AND type = ?", exchange, symbolType).Find(&symbols)
	if result.Error != nil {
		return nil, result.Error
	}

	return symbols, nil
}

type DataComparisonResult struct {
	Exchange      string
	Type          string
	APICount      int
	DBCount       int
	NewInAPI      []string
	MissingInAPI  []string
	CommonSymbols []string
}

func (p *Processor) CompareAPIWithDatabase(apiSymbols []models.Symbol, exchange, symbolType string) (*DataComparisonResult, error) {
	log.Printf("=== 开始对比%s %s数据 ===", exchange, symbolType)
	log.Printf("API获取到 %d 个交易对", len(apiSymbols))

	dbSymbols, err := p.getSymbolsByExchangeAndType(exchange, symbolType)
	if err != nil {
		log.Printf("获取数据库中%s %s交易对失败: %v", exchange, symbolType, err)
		return nil, err
	}

	log.Printf("数据库中有 %d 个%s %s交易对", len(dbSymbols), exchange, symbolType)

	apiSet := make(map[string]bool)
	dbSet := make(map[string]bool)

	for _, symbol := range apiSymbols {
		apiSet[symbol.Symbol] = true
	}

	for _, symbol := range dbSymbols {
		dbSet[symbol.Symbol] = true
	}

	var newInAPI []string
	var missingInAPI []string
	var common []string

	for symbol := range apiSet {
		if !dbSet[symbol] {
			newInAPI = append(newInAPI, symbol)
		} else {
			common = append(common, symbol)
		}
	}

	for symbol := range dbSet {
		if !apiSet[symbol] {
			missingInAPI = append(missingInAPI, symbol)
		}
	}

	sort.Strings(newInAPI)
	sort.Strings(missingInAPI)
	sort.Strings(common)

	result := &DataComparisonResult{
		Exchange:      exchange,
		Type:          symbolType,
		APICount:      len(apiSymbols),
		DBCount:       len(dbSymbols),
		NewInAPI:      newInAPI,
		MissingInAPI:  missingInAPI,
		CommonSymbols: common,
	}

	p.logDataComparisonResults(result)

	return result, nil
}

func (p *Processor) logDataComparisonResults(result *DataComparisonResult) {
	log.Printf("\n=== %s %s 数据对比详细结果 ===", result.Exchange, result.Type)
	log.Printf("API数据: %d 个交易对", result.APICount)
	log.Printf("数据库数据: %d 个交易对", result.DBCount)
	log.Printf("新增交易对: %d 个", len(result.NewInAPI))
	log.Printf("下架交易对: %d 个", len(result.MissingInAPI))
	log.Printf("保持不变: %d 个", len(result.CommonSymbols))

	if len(result.NewInAPI) > 0 {
		log.Printf("\n=== API中新增的交易对 (%d个) ===", len(result.NewInAPI))
		for i, symbol := range result.NewInAPI {
			if i < 15 {
				log.Printf("  [新增] %s", symbol)
			} else if i == 15 {
				log.Printf("  ... 还有 %d 个新增交易对", len(result.NewInAPI)-15)
				break
			}
		}
	}

	if len(result.MissingInAPI) > 0 {
		log.Printf("\n=== API中消失的交易对 (可能已下架) (%d个) ===", len(result.MissingInAPI))
		for i, symbol := range result.MissingInAPI {
			if i < 15 {
				log.Printf("  [消失] %s", symbol)
			} else if i == 15 {
				log.Printf("  ... 还有 %d 个消失的交易对", len(result.MissingInAPI)-15)
				break
			}
		}
	}

	changeRate := 0.0
	if result.DBCount > 0 {
		changeRate = float64(len(result.NewInAPI)+len(result.MissingInAPI)) / float64(result.DBCount) * 100
	}
	log.Printf("\n数据变化率: %.2f%% (变化数量/原数据库数量)", changeRate)

	log.Println("========================================")
}
