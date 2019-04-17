package controller

import (
	pb "../../common/rpc"
	"golang.org/x/net/context"
)

type RpcController struct {
}

func (controller *RpcController) execute(ctx context.Context, request *pb.GrpcRequestMsgInfo) (*pb.GrpcResponseMsgInfo, error) {
	requestBody := request.RequestBody
	serviceName := request.ServiceName


	ctx.()
	return &pb.GrpcResponseMsgInfo{ResponseBody: ""}, nil
}