package error

import "errors"

var (
	ErrOpQueryParams              = errors.New("transaction: error required Operation parameter")
	ErrOpParams                   = errors.New("transaction: error Operation type must be 'S' or 'B'")
	ErrCoinQueryParams            = errors.New("transaction: error required Coin parameter")
	ErrPriceQueryParams           = errors.New("transaction: error required Price parameter and must be greater than 0")
	ErrQtyQueryParams             = errors.New("transaction: error required Quantity parameter  and must be greater than 0")
	ErrInsufficientQtyQueryParams = errors.New("transaction: error Insufficient quantity for sale")

	ErrParamsBindValidate = errors.New("transaction: error fail to bind values")
)
