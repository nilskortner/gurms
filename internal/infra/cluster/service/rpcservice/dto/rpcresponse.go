package dto

type RpcResponse[T comparable] struct {
	RequestId int
	Result    T
	Rpcerror  error
}
