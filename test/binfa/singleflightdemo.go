package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"log"
	"strconv"
	"sync"
	"time"
)

var errorNotExist = errors.New("data not exist")

var g singleflight.Group

func main() {
	Single()
}

func singTimeout()  {
	ctx,cancel := context.WithTimeout(context.Background(),500*time.Millisecond)
	defer cancel()
	s, e := singleflightTimeout(ctx, &g, 15)
	if e!=nil{
		panic(e)
	}
	fmt.Println(s)
}

func singleflightTimeout(ctx context.Context, sg *singleflight.Group, id int) (string, error) {
	result := sg.DoChan(strconv.Itoa(id), func() (interface{}, error) {
		// 模拟出现问题，hang 住
		select {}
		return "",errors.New("error")
	})

	select {
	case r := <-result:
		return r.Val.(string), r.Err
	case <-ctx.Done():
		fmt.Println("timeout")
		return "", ctx.Err()
	}
}

func Single()  {
	var wg sync.WaitGroup
	wg.Add(10)

	//模拟10个并发
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("sunyuqi")
			if err != nil {
				log.Print(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}

//获取数据
func getData(key string) (string, error) {
	data, err := getDataFromCache(key)
	if errors.Is(err,errorNotExist) {
		//模拟从db中获取数据
		v, err, _ := g.Do(key, func() (interface{}, error) {
			return getDataFromDB(key)
			//set cache
		})
		if err != nil {
			log.Println(err)
			return "", err
		}
		//TOOD: set cache
		data = v.(string)
	} else if err != nil { //TOOD: set cache
		return "", err
	}
	return data, nil
}

//模拟从cache中获取值，cache中无该值
func getDataFromCache(key string) (string, error) {

	return "", errors.Wrap(errorNotExist,"user not fount")
}

//模拟从数据库中获取值
func getDataFromDB(key string) (string, error) {
	log.Printf("get %s from database", key)
	time.Sleep(time.Second)
	return key, nil
}