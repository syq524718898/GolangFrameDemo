package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"GolangFrameDemo/api"
	"GolangFrameDemo/internal/dao"
	"GolangFrameDemo/internal/service"
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
		fmt.Println(md.Get("name")[0])
		fmt.Println(md.Get("age")[0])
	}
	service := service.GetService()
	student, e := service.GetStudentById(int(a.GetId()))
	if e!=nil{
		return nil,e
	}
	rep := &api.UserReply{Id:a.Id,Name:student.S_name,Age:int32(student.D_id)}
	return rep, nil
}

func main() {
	d,con,err := dao.Init("configs/db.yml","configs/redis.yml")
	if err!=nil{
		panic(err)
	}
	dao, err := dao.New(d, con)
	if err!=nil {
		panic(err)
	}
	_, cfg, err := service.New(dao)
	if err!=nil{
		panic(err)
	}
	defer cfg()
	// 监听本地端口
	listener, err := net.Listen(UserMethod, UserAddress)
	if err != nil {
		panic(err)
		return
	}
	ser := grpc.NewServer()                       // 创建GRPC
	api.RegisterUserServiceServer(ser, &userserver{}) // 在GRPC服务端注册服务

	reflection.Register(ser) // 在GRPC服务器注册服务器反射服务
	// Serve方法接收监听的端口,每到一个连接创建一个ServerTransport和server的grroutine
	// 这个goroutine读取GRPC请求,调用已注册的处理程序进行响应
	err = ser.Serve(listener)
	if err != nil {
		panic(err)
		return
	}
}
