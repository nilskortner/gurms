package dto

type RpcResponseWrap interface {
	GetRpcError() error
}

type RpcResponse struct {
	RequestId int64
	Result    any
	RpcError  error
}

func (r *RpcResponse) GetRpcError() error {
	return r.RpcError
}
