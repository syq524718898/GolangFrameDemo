package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"math/rand"
	"net"
	"time"
	"GolangFrameDemo/api"
)
const (
	// Address 监听地址
	UserAddress string = "localhost:8080"
	// Method 通信方法
	UserMethod string = "tcp"
)


type userserver struct{
	api.UnimplementedUserServiceServer
}

// 接收client端的请求,函数名需保持一致
// ctx参数必传
// 参数二为自定义的参数,需从pb文件导入,因此pb文件必须可导入,文件放哪里随意
// 返回值同参数二,为pb文件的返回结构体指针
func (s *userserver) GetUser(ctx context.Context, a *api.UserRequest) (*api.UserReply, error) {
	// 逻辑写在这里
	fmt.Println("进入server逻辑")
	md, ok := metadata.FromIncomingContext(ctx)
	if ok{
		fmt.Println(md.Get("name"))
		fmt.Println(md.Get("age")[0])
	}
	// 模拟失败
	intn := rand.Intn(5)
	fmt.Printf("随机数：%d\n",intn)
	if intn<10 {
		return nil,status.Error(codes.NotFound,"user not fount")
	}
	// 模拟耗时
	time.Sleep(time.Duration(intn)*time.Millisecond)
	rep := &api.UserReply{Id:a.GetId(),Name:"sunyuqi",Age:18}
	return rep, nil
}



func main() {
	// 监听本地端口
	listener, err := net.Listen(UserMethod, UserAddress)
	if err != nil {
		panic(err)
		return
	}
	s := grpc.NewServer()                       // 创建GRPC
	api.RegisterUserServiceServer(s, &userserver{}) // 在GRPC服务端注册服务

	reflection.Register(s) // 在GRPC服务器注册服务器反射服务
	// Serve方法接收监听的端口,每到一个连接创建一个ServerTransport和server的grroutine
	// 这个goroutine读取GRPC请求,调用已注册的处理程序进行响应
	rand.Seed(time.Now().Unix())
	err = s.Serve(listener)
	if err != nil {
		panic(err)
		return
	}
}