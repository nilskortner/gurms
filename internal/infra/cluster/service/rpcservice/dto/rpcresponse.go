package dto

type RpcResponse struct {
	RequestId int
	Result    any
	Rpcerror  error
}
