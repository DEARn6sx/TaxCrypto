package controllers

import (
	"context"
	"extaxcrypto/db"
	errorTransaction "extaxcrypto/error/transaction"
	"extaxcrypto/models"
	"extaxcrypto/payloads"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func AddTransaction(c echo.Context) error {
	var req payloads.TransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorCaseTransaction(errorTransaction.ErrParamsBindValidate))
	}

	// Check if any of the required fields are empty
	if req.Op == "" {
		logrus.Errorln("Operation type is required but was not provided.")
		return c.JSON(http.StatusBadRequest, errorCaseTransaction(errorTransaction.ErrOpQueryParams))
	}
	if req.Op != "S" && req.Op != "B" {
		logrus.Errorln("Operation type must be 'S' or 'B'.")
		return c.JSON(http.StatusBadRequest, errorCaseTransaction(errorTransaction.ErrOpParams))
	}
	if req.Coin == "" {
		logrus.Errorln("Coin symbol is required but was not provided.")
		return c.JSON(http.StatusBadRequest, errorCaseTransaction(errorTransaction.ErrCoinQueryParams))
	}
	if req.Price <= 0 {
		logrus.Errorln("Price is required  must be greater than zero.")
		return c.JSON(http.StatusBadRequest, errorCaseTransaction(errorTransaction.ErrPriceQueryParams))
	}
	if req.Qty <= 0 {
		logrus.Errorln("Quantity is required must be greater than zero.")
		return c.JSON(http.StatusBadRequest, errorCaseTransaction(errorTransaction.ErrQtyQueryParams))
	}

	// Calculate TotalPrice
	totalPrice := req.Price * req.Qty

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
