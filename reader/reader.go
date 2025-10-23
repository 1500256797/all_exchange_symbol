package reader

import (
	"all_exchange_symbol/exchanges"
	"all_exchange_symbol/models"
	"log"
	"sync"
)

type Reader struct {
	exchanges []exchanges.ExchangeInterface
}

func NewReader() *Reader {
	return &Reader{
		exchanges: []exchanges.ExchangeInterface{
			exchanges.NewBinance(),
			exchanges.NewOKX(),
			exchanges.NewGate(),
			exchanges.NewBitget(),
			exchanges.NewBybit(),
		},
	}
}

func (r *Reader) FetchAllSymbols() ([]models.Symbol, error) {
	var allSymbols []models.Symbol
	var mu sync.Mutex
	var wg sync.WaitGroup
	errorChan := make(chan error, len(r.exchanges)*2)

	for _, exchange := range r.exchanges {
		wg.Add(2)

		go func(ex exchanges.ExchangeInterface) {
			defer wg.Done()
			symbols, err := ex.FetchSpotSymbols()
			if err != nil {
				log.Printf("Error fetching spot symbols from %s: %v", ex.GetName(), err)
				errorChan <- err
				return
			}

			mu.Lock()
			allSymbols = append(allSymbols, symbols...)
			mu.Unlock()

			log.Printf("Successfully fetched %d spot symbols from %s", len(symbols), ex.GetName())
		}(exchange)

		go func(ex exchanges.ExchangeInterface) {
			defer wg.Done()
			symbols, err := ex.FetchFuturesSymbols()
			if err != nil {
				log.Printf("Error fetching futures symbols from %s: %v", ex.GetName(), err)
				errorChan <- err
				return
			}

			mu.Lock()
			allSymbols = append(allSymbols, symbols...)
			mu.Unlock()

			log.Printf("Successfully fetched %d futures symbols from %s", len(symbols), ex.GetName())
		}(exchange)
	}

	wg.Wait()
	close(errorChan)

	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		log.Printf("Encountered %d errors during fetching, but continuing with available data", len(errors))
	}

	log.Printf("Total symbols fetched: %d", len(allSymbols))
	return allSymbols, nil
}

func (r *Reader) FetchSymbolsByExchange(exchangeName string) ([]models.Symbol, error) {
	for _, exchange := range r.exchanges {
		if exchange.GetName() == exchangeName {
			var allSymbols []models.Symbol

			spotSymbols, err := exchange.FetchSpotSymbols()
			if err != nil {
				log.Printf("Error fetching spot symbols from %s: %v", exchangeName, err)
			} else {
				allSymbols = append(allSymbols, spotSymbols...)
			}

			futuresSymbols, err := exchange.FetchFuturesSymbols()
			if err != nil {
				log.Printf("Error fetching futures symbols from %s: %v", exchangeName, err)
			} else {
				allSymbols = append(allSymbols, futuresSymbols...)
			}

			return allSymbols, nil
		}
	}

	return nil, nil
}
