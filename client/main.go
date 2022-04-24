package main

import (
	"context"
	"errors"
	"fmt"
	xerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
	"GolangFrameDemo/api"
)

const (
	// Address server端地址
	Address string = "localhost:8080"
)

func main() {
	reply, e := CallRpcTest()
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

func CallRpcTest()(*api.HelloReply,error)  {
	// 连接服务器
	conn, err := grpc.Dial(Address,grpc.WithInsecure())
	if err != nil {
		return nil,xerrors.Wrap(errors.New("connect err"),"RPC连接失败")
	}
	defer conn.Close()

	// 连接GRPC
	c := api.NewGreeterClient(conn)
	// 创建要发送的结构体
	req := api.HelloRequest{Name:"张三"}
	// 携带元数据
	md := metadata.Pairs("name", "孙渝其", "age", "18")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 超时控制, 1秒超时
	ctx,cancel := context.WithTimeout(ctx,400*time.Millisecond)
	defer cancel()

	var (
		r *api.HelloReply
		rchan = make(chan *api.HelloReply)
	)
	go func() {
		defer close(rchan)
		// 调用server的注册方法
		reply, inierr := c.SayHello(ctx, &req)
		if inierr != nil {
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