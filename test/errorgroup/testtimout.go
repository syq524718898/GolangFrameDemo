
package main

import (
	"context"
	"errors"
	"fmt"
	xerrors "github.com/pkg/errors"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	RpcTest()
}

func testChan()  {
	c := make(chan int)
	go func() {
		time.Sleep(2*time.Second)
		fmt.Println(<-c)
		fmt.Println("chan 关闭了")
	}()
	close(c)
	time.Sleep(2*time.Second)
	fmt.Println("程序结束")
}

func testWg()  {
	rand.Seed(time.Now().Unix())
	var wg = new(sync.WaitGroup)
	wg.Add(2)
	go testWGFunc(wg)
	go testWGFunc(wg)
	wg.Wait()
	fmt.Println("结束")

}
func testWGFunc(wg *sync.WaitGroup)  {
	defer wg.Done()
	i := rand.Intn(5)+1
	time.Sleep(time.Second*time.Duration(i))
	fmt.Printf("完成任务,耗时%d秒\n",i)
}



type User struct {
	name string
	age string
}

func RpcTest()  {
	user, e := CallRpc(context.Background())
	time.Sleep(4*time.Second)
	if e!=nil{
		fmt.Printf("original error:%T %v\n",xerrors.Cause(e),xerrors.Cause(e))
		fmt.Printf("stack trace:\n%+v\n",e)
		os.Exit(1)
	}
	fmt.Println(user)
}

// 模拟调用
func CallRpc(ctx context.Context)(*User,error)  {
	// 超时控制
	ctx,cancel := context.WithTimeout(ctx,1*time.Second)
	defer cancel()
	var (
		uchan = make(chan *User)
		err error
		u *User
	)

	go func() {
		defer close(uchan)
		// 调用server的注册方法
		ctx = context.WithValue(ctx, "name", "孙渝其")
		reply, inlineerr := GetUser(ctx)
		if inlineerr != nil {
			err = xerrors.WithMessage(inlineerr,"rpc调用返回错误")
			return
		}
		uchan<-reply
	}()

	select {
	case <-ctx.Done():
		fmt.Println("超时啦")
		err = xerrors.Wrap(errors.New("timeout"),"超时")
	case u=<-uchan:
		fmt.Println("rpc success")
	}
	return u,err
}

// 服务
func GetUser(ctx context.Context)(*User,error) {
	if rand.Intn(5)<2{
		time.Sleep(500*time.Millisecond)
		u := &User{name:ctx.Value("name").(string),age:"18"}
		return u,nil
	}else {
		return nil,xerrors.Wrap(errors.New("sql error"),"查询数据库错误")
	}
}