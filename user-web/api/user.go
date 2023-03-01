package api

import (
	"context"
	"fmt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"go_mxshop_web/user-web/forms"
	"go_mxshop_web/user-web/global"
	"go_mxshop_web/user-web/middlewares"
	"go_mxshop_web/user-web/models"
	"go_mxshop_web/user-web/proto"
	"go_mxshop_web/user-web/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"
)

// 将grpc错误状态码转换为http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err == nil {
		return
	}

	if e, ok := status.FromError(err); ok {
		zap.S().Errorf("HandleGrpcErrorToHttp: %d, %s", e.Code(), e.Message())
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

// 删除表单异常信息中部分
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}

	return rsp
}

func GetUserList(ctx *gin.Context) {

	ip := global.ServerConfig.UserSrv.Host
	port := global.ServerConfig.UserSrv.Port
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

func PassWordLogin(ctx *gin.Context) {
	request := forms.PassWordLoginForm{}
	if err := ctx.ShouldBind(&request); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": removeTopStruct(errs.Translate(global.Trans)),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}

	// ip := global.ServerConfig.UserSrv.Host
	// port := global.ServerConfig.UserSrv.Port
	// conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("consul://192.168.168.180:8500/user-srv?wait=14s&tag=manual",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("【GetUserList】 连接 【用户服务】失败", "msg", err.Error())
	}
	client := proto.NewUserClient(conn)
	rsp, err := client.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: request.Mobile,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	if rsp.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "用户信息不存在"})
		return
	}

	checkRsp, err := client.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
		Password:          request.PassWord,
		EncryptedPassword: rsp.Password,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	if !checkRsp.Success {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "密码输入有误"})
		return
	}

	claims := models.CustomClaims{
		ID:       uint(rsp.Id),
		NickName: rsp.NickName,
		StandardClaims: jwt2.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24,
			Issuer:    "imooc",
		},
	}
	jwt := middlewares.NewJWT()
	token, err := jwt.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "token生成失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":        uint(rsp.Id),
		"nickName":  rsp.NickName,
		"token":     token,
		"expiresAt": (time.Now().Unix() + 60*60*24) * 1000,
	})
}

func CreateUser(ctx *gin.Context) {
	userRequest := forms.CreateUserForm{}
	if err := ctx.ShouldBind(&userRequest); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": removeTopStruct(errs.Translate(global.Trans)),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}

	//ip := global.ServerConfig.UserSrv.Host
	//port := global.ServerConfig.UserSrv.Port
	conn, err := grpc.Dial(
		"consul://192.168.168.180:8500/user-srv?wait=14s&tag=manual",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("【GetUserList】 连接 【用户服务】失败", "msg", err.Error())
	}
	client := proto.NewUserClient(conn)
	rsp, err := client.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   userRequest.Mobile,
		NickName: userRequest.NickName,
		Password: userRequest.PassWord,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": rsp.Id})
}
