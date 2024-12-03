package dto

type RpcResponseWrap interface {
	GetRpcError() error
}

type RpcResponse[T comparable] struct {
	RequestId int64
	Result    T
	RpcError  error
}

func (r *RpcResponse[T]) GetRpcError() error {
	return r.RpcError
}
