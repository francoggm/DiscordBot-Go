package stockmarket

import (
	"discord-bot/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Ticker struct {
	Symbol      string `json:"1. symbol"`
	Name        string `json:"2. name"`
	Type        string `json:"3. type"`
	Region      string `json:"4. region"`
	OpenMarket  string `json:"5. marketOpen"`
	CloseMarket string `json:"6. marketClose"`
	TimeZone    string `json:"7. timezone"`
	Currency    string `json:"8. currency"`
	Score       string `json:"9. matchScore"`
}

type FindTicker struct {
	Tickers []Ticker `json:"bestMatches"`
}

const (
	baseURL            = "https://www.alphavantage.co/query"
	findTickerFunction = "SYMBOL_SEARCH"
)

func FindTickers(keywords []string) []Ticker {
	cfg := config.GetConfig()
	var tickers []Ticker

	for _, keyword := range keywords {
		if keyword != "" {
			URL := fmt.Sprintf("%s?function=%s&keywords=%s&apikey=%s", baseURL, findTickerFunction, keyword, cfg.StockAPIKey)

			res, err := http.Get(URL)
			if err != nil {
				continue
			}

			if res.StatusCode != http.StatusOK {
				continue
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				continue
			}

			var tks FindTicker

			err = json.Unmarshal(body, &tks)
			if err != nil {
				continue
			}

			tickers = append(tickers, tks.Tickers...)
		}
	}

	return tickers
}
