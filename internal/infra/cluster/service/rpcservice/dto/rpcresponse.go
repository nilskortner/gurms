package dto

type RpcResponse[T comparable] struct {
	RequestId int64
	Result    T
	Rpcerror  error
}
