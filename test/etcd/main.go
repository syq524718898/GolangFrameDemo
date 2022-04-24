package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"

	"strconv"
	"time"
)

func main() {
	endpoints := []string{
		"127.0.0.1:20000",
		"127.0.0.1:20002",
		"127.0.0.1:20004",
	}
	client3, err := clientv3.New(
		clientv3.Config{
			Endpoints: endpoints,
		},
	)

	if err != nil {
		panic(err.Error())
	}


	// etcdctl put  写入
	response, err := client3.Put(context.TODO(), "name", "zhangsan")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(response)

	// etcdctl put  写入
	response4, err := client3.Put(context.TODO(), "web4", "web4444444")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(response4)

	// etcdctl get 获取
	fmt.Println("etcdctl get 获取")
	getResponse, err := client3.Get(context.TODO(), "name")
	if err != nil {
		panic(err.Error())
	}
	for _, kv := range getResponse.Kvs {
		fmt.Println(kv)
		fmt.Println(string(kv.Key), string(kv.Value))
	}

	// etcdctl get --prefix 获取带前缀的key的值
	fmt.Println("etcdctl get --prefix 获取带前缀的key的值")
	getAll, err := client3.Get(context.TODO(), "web", clientv3.WithPrefix())
	if err != nil {
		panic(err.Error())
	}
	for _, kv1 := range getAll.Kvs {
		fmt.Println(kv1)
		fmt.Println(string(kv1.Key), string(kv1.Value))
	}

	// etcdctl del 刪除元素
	deleteResponse, err := client3.Delete(context.TODO(), "name")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(deleteResponse)

	// etcdctl watch 监控一个key的变化 程序会卡在这不动 一直监控着watch监控的key值的变化 通过channel管道进行的监控 所以无时无刻不在监控 所以呢主进程不会死掉的！
	//watch := client3.Watch(context.TODO(), "web")
	//for wc := range watch {
	//  for _, w := range wc.Events {
	//      fmt.Println(string(w.Kv.Key), string(w.Kv.Value), w.Type.String())
	//  }
	//}

	// etcdctl lease grant 20 租赁时长是20s 20s之后就自动过期了！
	grantResponse, err := client3.Grant(context.TODO(), 20)
	if err != nil {
		panic(err.Error())
	}
	if _, err := client3.Put(context.TODO(), "wangxuancheng", "abcdefg", clientv3.WithLease(grantResponse.ID)); err != nil {
		panic(err.Error())
	}
	fmt.Println("ok")

	// 测试租赁是否成功
	i := 0
	for {
		ps, err := client3.Get(context.TODO(), "wangxuancheng")
		if err != nil {
			panic(err.Error())
		}
		for _, kv := range ps.Kvs {
			fmt.Println(string(kv.Key) + "----" + string(kv.Value))
		}
		i += 5
		fmt.Println(strconv.Itoa(i) + "s")
		time.Sleep(5 * time.Second)
	}

}

