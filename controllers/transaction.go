package controllers

import (
	"context"
	"extaxcrypto/db"
	"extaxcrypto/models"
	"extaxcrypto/payloads"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddTransaction(c echo.Context) error {
	var req payloads.TransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Calculate TotalPrice
	totalPrice := req.Price * req.Qty

	if req.Op == "S" {
		// Calculate the remaining quantity of the coin
		matchStage := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "coin", Value: req.Coin},
			}},
		}
		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$coin"},
				{Key: "totalQty", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$cond", Value: bson.D{
							{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$op", "B"}}}},
							{Key: "then", Value: "$qty"},
							{Key: "else", Value: bson.D{{Key: "$multiply", Value: bson.A{"$qty", -1}}}},
						}},
					}},
				}},
			}},
		}

		cursor, err := db.Collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer cursor.Close(context.TODO())

		var results []struct {
			TotalQty float64 `bson:"totalQty"`
		}
		if err := cursor.All(context.TODO(), &results); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		remainingQty := 0.0
		if len(results) > 0 {
			remainingQty = results[0].TotalQty
		}

		if remainingQty < req.Qty {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error":         "Insufficient quantity for sale",
				"remaining_qty": remainingQty,
			})
		}
	}

	// Insert the transaction
	_, err := db.Collection.InsertOne(context.TODO(), models.Transaction{
		Op:         req.Op,
		Coin:       req.Coin,
		Price:      req.Price,
		Qty:        req.Qty,
		TotalPrice: totalPrice,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := payloads.TransactionResponse{
		Op:         req.Op,
		Coin:       req.Coin,
		Price:      req.Price,
		Qty:        req.Qty,
		TotalPrice: totalPrice,
	}

	return c.JSON(http.StatusOK, response)
}

func CalculateAveragePrice(c echo.Context) error {
	cursor, err := db.Collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer cursor.Close(context.TODO())

	var transactions []models.Transaction
	if err = cursor.All(context.TODO(), &transactions); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	averagePricesBuy := make(map[string]payloads.AveragePriceResponse)
	averagePricesSell := make(map[string]payloads.AveragePriceResponse)
	remainingQty := make(map[string]float64)

	for _, txn := range transactions {
		coin := txn.Coin
		switch txn.Op {
		case "B":
			// Process buy transactions
			avgPriceData, exists := averagePricesBuy[coin]
			if !exists {
				avgPriceData = payloads.AveragePriceResponse{Coin: coin}
			}
			avgPriceData.Qty += txn.Qty
			avgPriceData.TotalPrice += txn.TotalPrice
			avgPriceData.AveragePrice = avgPriceData.TotalPrice / avgPriceData.Qty
			averagePricesBuy[coin] = avgPriceData

			// Update remaining quantity
			remainingQty[coin] += txn.Qty

		case "S":
			// Process sell transactions
			avgPriceData, exists := averagePricesSell[coin]
			if !exists {
				avgPriceData = payloads.AveragePriceResponse{Coin: coin}
			}
			avgPriceData.Qty += txn.Qty
			avgPriceData.TotalPrice += txn.TotalPrice
			avgPriceData.AveragePrice = avgPriceData.TotalPrice / avgPriceData.Qty
			averagePricesSell[coin] = avgPriceData

			// Update remaining quantity
			remainingQty[coin] -= txn.Qty
		}
	}

	response := payloads.AveragePriceDetailsResponse{
		AveragePricesBuy:  averagePricesBuy,
		AveragePricesSell: averagePricesSell,
		RemainingQty:      remainingQty,
	}

	return c.JSON(http.StatusOK, response)
}
