package errorStandard

import "errors"

var (
	ErrInternalServer    = errors.New("server: error the system")
	ErrRequiredParameter = errors.New("server: error required parameter")
	ErrMarshal           = errors.New("server: error json marshal")
	ErrUnmarshal         = errors.New("server: error json unmarshal")
	ErrQueryParams       = errors.New("server: error query params is required")
	ErrFuncNotAll        = errors.New("server: error func not all")
	ErrEncode            = errors.New("server: error encode")
	ErrDecode            = errors.New("server: error decode")
	ErrEncodeToken       = errors.New("server: error encode token")
	ErrDecodeToken       = errors.New("server: error decode token")
)
