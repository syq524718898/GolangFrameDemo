package dao

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type RedisConf struct {
	Proto 		string	`yaml:"proto"`
	Addr 	string	`yaml:"addr"`
}


func Redisconf(path string) (*RedisConf,error)  {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil,errors.New("file open err")
	}
	redisConf := &RedisConf{}
	err = yaml.Unmarshal(yamlFile, redisConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil,errors.New("yml ummarshal err")
	}
	fmt.Println(redisConf)
	return redisConf,nil
}

func RedisInit(path string)(con redis.Conn,err error)  {
	conf, err := Redisconf(path)
	if err!=nil{
		return
	}
	con, err = redis.Dial(conf.Proto, conf.Addr)
	if err != nil {
		return
	}
	return
}
