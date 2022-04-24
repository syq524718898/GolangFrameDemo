package main

import (
	"context"
	"errors"
	"fmt"
	xerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"time"
	"GolangFrameDemo/api"
)

const (
	// Address server端地址
	UserAddress string = "localhost:8080"
)


func main() {

	reply, e := CallUserRpc()
	if e!=nil{
		fmt.Printf("original error:%T %v\n",xerrors.Cause(e),xerrors.Cause(e))
		fmt.Printf("stack trace:\n%+v\n",e)
		os.Exit(1)
	}else {
		fmt.Println(reply)
	}
	time.Sleep(3*time.Second)
	fmt.Println("服务结束")
}

func CallUserRpc()(*api.UserReply,error)  {
	// 连接服务器
	conn, err := grpc.Dial(UserAddress,grpc.WithInsecure())
	if err != nil {
		return nil,xerrors.Wrap(errors.New("connect err"),"RPC连接失败")
	}
	defer conn.Close()

	// 连接GRPC
	c := api.NewUserServiceClient(conn)

	// 创建要发送的结构体
	req := api.UserRequest{Id:1}
	// 携带元数据
	md := metadata.Pairs("age", "18")
	md.Set("name","孙渝其","祝依杰")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 超时控制, 1秒超时
	ctx,cancel := context.WithTimeout(ctx,1000*time.Millisecond)
	defer cancel()

	var (
		r *api.UserReply
		rchan = make(chan *api.UserReply)
	)
	go func() {
		defer close(rchan)
		// 调用server的注册方法
		reply, inierr := c.GetUser(ctx, &req)
		if inierr != nil {
			if status.Code(inierr) == codes.NotFound{
				fmt.Println("=============")
				fmt.Println("not fount")
				fmt.Println("=============")
			}
			err = xerrors.Wrap(inierr,"RPC调用错误")
			return
		}
		rchan<-reply
	}()
	select {
	case <-ctx.Done():
		err = xerrors.Wrap(errors.New("timeout"),"RPC调用超时")
	case r=<-rchan:

	}
	return r,err
}