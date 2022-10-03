package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	//初始化 limiter 每秒1个令牌，令牌桶容量为1
	limiter :=  rate.NewLimiter(rate.Every(time.Millisecond*1000),1)


	//阻塞直到获取足够的令牌或者上下文取消
	ctx,_ := context.WithTimeout(context.Background(),time.Second*10)
	for {
		// 限流 每秒一次
		err := limiter.Wait(ctx)
		if err != nil {
			fmt.Println("error",err)
		}
		fmt.Println("hhh")
		time.Sleep(200*time.Millisecond)
	}
}
