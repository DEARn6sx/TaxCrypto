package errorExternal

import "errors"

var (
	ErrCallExternalChangeMeAPI     = errors.New("server: error call external api change me ")
	ErrTimeOutCallExternalChangeMe = errors.New("server: error call external api change me time out")
	ErrMarshalExternalChangeMe     = errors.New("server: error json marshal external api change me")
	ErrUnmarshalExternalChangeMe   = errors.New("server: error json unmarshal external api change me")
)
