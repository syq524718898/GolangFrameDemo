package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"GolangFrameDemo/api"
	"GolangFrameDemo/test/encoding"
)




type Student struct {
	Id int `json:"id"`
	S_name string `json:"s_name"`
	D_id int `json:"d_id"`
}

type RedisConf struct {
	Proto 		string	`yaml:"proto"`
	Addr 	string	`yaml:"addr"`
}

type MysqlConf struct {
	DriverName 		string	`yaml:"driverName"`
	DataSourceName 	string	`yaml:"dataSourceName"`
}

func ReadMysqlYml(path string) (*MysqlConf,error)  {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil,errors.New("file open err")
	}
	mysqlConf := &MysqlConf{}
	err = yaml.Unmarshal(yamlFile, mysqlConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil,errors.New("yml ummarshal err")
	}
	fmt.Println(mysqlConf)
	return mysqlConf,nil
}

func main() {
	//WriteRedis(&u,encoding.JSON)
	//TestRedis()
	TestMysql()
}

func TestCode()  {

	u := &api.UserReply{
		Id:*proto.Int64(1),
		Name:*proto.String("孙渝其"),
		Age:*proto.Int32(18),
	}
	t := encoding.JSON
	bytes, e := encoding.Marsh(u, t)
	if e!=nil{
		panic(e)
	}
	reply, e := encoding.UnMarsh(bytes, t)
	if e!=nil{
		panic(e)
	}
	fmt.Println(reply)
	fmt.Println(reply.Id)
	fmt.Println(reply.Age)
	fmt.Println(reply.Name)
}

func TestMysql()  {
	conf,e := ReadMysqlYml("configs/db.yml")
	if e!=nil {
		panic(e)
	}
	db, e := sql.Open(conf.DriverName,conf.DataSourceName)
	defer db.Close()
	if e!=nil{
		fmt.Println("数据库连接失败")
		panic(e)
	}

	rows,e:=db.Query("SELECT * FROM student where  id = 10000")

	if e!=nil{
		fmt.Println("出错")
		panic(e)
	}
	var s Student

	for rows.Next(){
		rows.Scan(&s.Id,&s.S_name,&s.D_id,)
		fmt.Println(s)
	}
	rows.Close()

	//WriteRedis(&s)
}




func ReadRedisYml(path string) (*RedisConf,error)  {
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

func WriteRedis(stu *api.UserReply,t encoding.MarshType)  {
	// 通过go向redis写入数据和读取数据
	// 1.连接到redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis connect failed. err = ", err)
		return
	}
	defer conn.Close() //关闭连接

	var bytes []byte
	// 2.通过go向redis写入数据
	switch t {
	case encoding.PB:
		bytes, err = proto.Marshal(stu)
	case encoding.JSON:
		bytes, err = json.Marshal(*stu)
	default:
		bytes, err = json.Marshal(*stu)
	}

	if err !=nil{
		panic(err)
	}
	_, err = conn.Do("Set", fmt.Sprintf("stu_%v",stu.Id), bytes)
	if err != nil {
		fmt.Println("Set err = ", err)
		return
	}
}

func TestRedis()  {
	// 通过go向redis写入数据和读取数据
	conf, e := ReadRedisYml("configs/redis.yml")
	if e!=nil{
		panic(e)
	}
	// 1.连接到redis
	conn, err := redis.Dial(conf.Proto, conf.Addr)
	if err != nil {
		fmt.Println("redis connect failed. err = ", err)
		return
	}
	defer conn.Close() //关闭连接

	// 从redis中读取数据
	// 由于返回的r是一个interface{}，故要转成字节数组
	reply, e := conn.Do("Get", "stu_15")
	if e!=nil{
		fmt.Println("err = ", e)
		return
	}
	r, err := redis.Bytes(reply,e)
	if err != nil {
		fmt.Println("Get err = ", err)
		return
	}
	var stu api.UserReply
	//_ = json.Unmarshal(r,&stu)
	stu,_ = encoding.UnMarsh(r,encoding.PB)
	fmt.Println(stu)
}
