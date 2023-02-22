package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_mxshop_web/user-web/proto"
	"go_mxshop_web/user-web/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

// 将grpc错误状态码转换为http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err == nil {
		return
	}

	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"msg": "其他错误"})
		}
	}
}

func GetUserList(ctx *gin.Context) {
	ip := "127.0.0.1"
	port := 50051
	zap.S().Info("获取用户列表页")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("【GetUserList】 连接 【用户服务】失败", "msg", err.Error())
	}
	client := proto.NewUserClient(conn)

	rsp, err := client.GetUserList(context.Background(), &proto.PageInfo{Pn: 0, PSize: 1})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]response.UserResponse, 0, len(rsp.Data))
	for _, value := range rsp.Data {
		result = append(result, response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: response.JsonTime(time.Unix(int64(value.BrithDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		})
	}

	ctx.JSON(http.StatusOK, result)
}
