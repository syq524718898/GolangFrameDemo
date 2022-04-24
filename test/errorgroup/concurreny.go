package main

import (
	"context"
	"errors"
	"fmt"
	xerrors "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var (
	rpcnum = 3
	timeout = 900 // ms
)

func main() {

	rand.Seed(time.Now().Unix())
	ErrGroupAndRpc()
	time.Sleep(time.Second)
}


type Pojo struct {
	name string
	age string
}



// 模拟并发发送三个rpc请求,合并返回
func ErrGroupAndRpc()  {

	var g errgroup.Group
	// 超时控制 950ms超时
	ctx,cancal := context.WithTimeout(context.Background(),time.Duration(timeout)*time.Millisecond)
	// 通道保存结果
	uchans := make(chan *Pojo,rpcnum)
	// 保存结果
	users := make([]*Pojo,rpcnum)
	defer cancal()
	defer close(uchans)
	for i:=0;i<rpcnum;i++{
		// 并发发送三个请求
		g.Go(func() error {
			uchan := make(chan *Pojo)
			var e error
			var user *Pojo
			Go(func() {
				defer close(uchan)
				user, e = GetPojo(ctx)
				if e!=nil{
					return
				}
				uchan<-user
			})
			select {
			case u:=<-uchan:
				uchans<-u
			case <-ctx.Done():
				return xerrors.Wrap(errors.New("timeout"),"超时")
			}
			return e
		})
	}
	if err := g.Wait();err!=nil{
		fmt.Printf("original error:%T %v\n",xerrors.Cause(err),xerrors.Cause(err))
		fmt.Printf("stack trace:\n%+v\n",err)
		return
	}
	for i:=0;i<rpcnum;i++{
		users[i] = <-uchans
		fmt.Println(users[i])
	}
}


// 服务
func GetPojo(ctx context.Context)(*Pojo,error) {
	intn := rand.Intn(100)
	if intn<99{
		t := rand.Intn(10)
		time.Sleep(time.Millisecond*time.Duration(t*100))
		u := &Pojo{name:fmt.Sprintf("sunyuqi_%d",intn),age:"18"}
		return u,nil
	}else {
		return nil,xerrors.Wrap(errors.New("sql error"),"查询数据库错误")
	}
}



func ErrGroupHttp() {
	var g errgroup.Group
	var urls = []string{
		"https://blog.csdn.net/yzf279533105/article/details/97039688",
		"http://www.baidu.com/",
		"https://www.bilibili.com/",
	}
	for _, url := range urls {
		url := url // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			// Fetch the URL.
			resp, err := http.Get(url)
			if err == nil {
				defer resp.Body.Close()
				bytes, _ := ioutil.ReadAll(resp.Body)
				fmt.Printf("%s\n",bytes)
			}
			return err
		})
	}
	// Wait for all HTTP fetches to complete.
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}
}



func Go(x func())  {
	go func() {
		defer func() {
			if err := recover();err!=nil{
				fmt.Printf("error happend....: %s",err)
			}
		}()
		x()
	}()
}