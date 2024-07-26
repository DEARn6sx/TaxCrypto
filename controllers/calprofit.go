package controllers

import (
	"context"
	"extaxcrypto/db"
	"extaxcrypto/models"
	"extaxcrypto/payloads"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func CalculateProfit(c echo.Context) error {
	cursor, err := db.Collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var transactions []models.Transaction
	if err = cursor.All(context.TODO(), &transactions); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	totalProfit, profitByCoin, detailedProfitLoss := CalculateProfitLoss(transactions)

	response := payloads.ProfitDetailsResponse{
		TotalProfit:        totalProfit,
		ProfitByCoin:       profitByCoin,
		DetailedProfitLoss: detailedProfitLoss,
	}

	return c.JSON(http.StatusOK, response)
}

func CalculateProfitLoss(transactions []models.Transaction) (float64, map[string]float64, map[string][]payloads.CoinProfitLoss) {
	buyOrders := make(map[string][]models.Transaction)               // Holds buy orders for each coin
	profitByCoin := make(map[string]float64)                         // Holds profit for each coin
	detailedProfitLoss := make(map[string][]payloads.CoinProfitLoss) // Detailed profit/loss for each coin
	var totalProfit float64

	for _, txn := range transactions {
		switch txn.Op {
		case "B":
			// Add buy transaction to the list for this coin
			buyOrders[txn.Coin] = append(buyOrders[txn.Coin], txn)
		case "S":
			// Process sell transaction
			if len(buyOrders[txn.Coin]) == 0 {
				fmt.Printf("No buy orders available for %s, skipping...\n", txn.Coin)
				continue
			}

			sellQty := txn.Qty
			sellPrice := txn.Price
			var profitLoss float64

			for len(buyOrders[txn.Coin]) > 0 && sellQty > 0 {
				buyTxn := &buyOrders[txn.Coin][0]
				if buyTxn.Qty > sellQty {
					// Partial sale from the current buy order
					profit := (sellPrice - buyTxn.Price) * sellQty
					profitLoss += profit
					buyTxn.Qty -= sellQty

					// Append to detailed profit/loss
					detailedProfitLoss[txn.Coin] = append(detailedProfitLoss[txn.Coin], payloads.CoinProfitLoss{
						BuyPrice:  buyTxn.Price,
						SellPrice: sellPrice,
						Qty:       sellQty,
						Profit:    profit,
					})
					sellQty = 0
				} else {
					// Fully use up the current buy order
					profit := (sellPrice - buyTxn.Price) * buyTxn.Qty
					profitLoss += profit
					sellQty -= buyTxn.Qty

					// Append to detailed profit/loss
					detailedProfitLoss[txn.Coin] = append(detailedProfitLoss[txn.Coin], payloads.CoinProfitLoss{
						BuyPrice:  buyTxn.Price,
						SellPrice: sellPrice,
						Qty:       buyTxn.Qty,
						Profit:    profit,
					})
					buyOrders[txn.Coin] = buyOrders[txn.Coin][1:] // Remove used buy order
				}
			}

			if sellQty > 0 {
				fmt.Printf("Not enough buy orders to fulfill sell for %s, skipping...\n", txn.Coin)
			}

			// Update total and per coin profit
			totalProfit += profitLoss
			profitByCoin[txn.Coin] += profitLoss
		}
	}

	return totalProfit, profitByCoin, detailedProfitLoss
}
