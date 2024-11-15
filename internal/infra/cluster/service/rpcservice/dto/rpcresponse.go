package dto

type RpcResponse struct {
	requestId int
	result    any
	rpcerror  error
}
