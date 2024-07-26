package payloads

type TransactionResponse struct {
	Op         string  `json:"op"`
	Coin       string  `json:"coin"`
	Price      float64 `json:"price"`
	Qty        float64 `json:"qty"`
	TotalPrice float64 `json:"total_price"`
}
type AveragePriceResponse struct {
	Coin         string  `json:"coin"`
	Qty          float64 `json:"qty"`
	TotalPrice   float64 `json:"total_price"`
	AveragePrice float64 `json:"average_price"`
}

type ProfitDetailsResponse struct {
	TotalProfit        float64                     `json:"total_profit"`
	ProfitByCoin       map[string]float64          `json:"profit_by_coin"`
	DetailedProfitLoss map[string][]CoinProfitLoss `json:"detailed_profit_loss"`
}

type CoinProfitLoss struct {
	BuyPrice  float64 `json:"buy_price"`
	SellPrice float64 `json:"sell_price"`
	Qty       float64 `json:"qty"`
	Profit    float64 `json:"profit"`
}

type AveragePriceDetailsResponse struct {
	AveragePricesBuy  map[string]AveragePriceResponse `json:"average_prices_buy"`
	AveragePricesSell map[string]AveragePriceResponse `json:"average_prices_sell"`
	RemainingQty      map[string]float64              `json:"remaining_qty"`
}
