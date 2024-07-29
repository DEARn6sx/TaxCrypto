package controllers

import (
	"errors"
	errorCase "extaxcrypto/error"
	errorTransaction "extaxcrypto/error/transaction"
)

func errorCaseTransaction(err error) *errorCase.ResponseError {
	switch {
	// TODO Fix Case 1-7

	case errors.Is(err, errorTransaction.ErrParamsBindValidate):
		res := errorCase.ResponseError{
			Status:  400,
			Code:    10100,
			LabelTH: "ไม่สามารถดำเนินการได้ เนื่องจากระบบไม่สามารถ Bind ค่าได้",
			LabelEN: "error fail to bind values",
		}
		return &res
	case errors.Is(err, errorTransaction.ErrOpQueryParams):
		res := errorCase.ResponseError{
			Status:  400,
			Code:    10101,
			LabelTH: "ไม่สามารถดำเนินการได้ เนื่องจากระบบต้องการค่า Operation ในพารามิเตอร์",
			LabelEN: "input Operation parameter is invalid or missing, please correct the input data.",
		}
		return &res
	case errors.Is(err, errorTransaction.ErrOpParams):
		res := errorCase.ResponseError{
			Status:  400,
			Code:    10102,
			LabelTH: "ไม่สามารถดำเนินการได้ เนื่องจากระบบต้องการค่า Operation เป็นค่า 'S' หรือ 'B' เท่านั้น",
			LabelEN: "input Operation type must be 'S' or 'B'.",
		}
		return &res
	case errors.Is(err, errorTransaction.ErrCoinQueryParams):
		res := errorCase.ResponseError{
			Status:  400,
			Code:    10103,
			LabelTH: "ไม่สามารถดำเนินการได้ เนื่องจากระบบต้องการค่า Coin ในพารามิเตอร์",
			LabelEN: "input Coin parameter is invalid or missing, please correct the input data.",
		}
		return &res
	case errors.Is(err, errorTransaction.ErrPriceQueryParams):
		res := errorCase.ResponseError{
			Status:  400,
			Code:    10104,
			LabelTH: "ไม่สามารถดำเนินการได้ เนื่องจากระบบต้องการค่า Price ในพารามิเตอร์",
			LabelEN: "input Price parameter is invalid or missing, please correct the input data.",
		}
		return &res
	case errors.Is(err, errorTransaction.ErrQtyQueryParams):
		res := errorCase.ResponseError{
			Status:  400,
			Code:    10105,
			LabelTH: "ไม่สามารถดำเนินการได้ เนื่องจากระบบต้องการค่า Quantity ในพารามิเตอร์",
			LabelEN: "input Quantity parameter is invalid or missing, please correct the input data",
		}
		return &res
	}
	// TODO Fix Case
	res := errorCase.ResponseError{
		Status:  500,
		Code:    900,
		LabelTH: "ระบบเกิดข้อผิดพลาด",
		LabelEN: "error the system",
	}
	return &res
}
