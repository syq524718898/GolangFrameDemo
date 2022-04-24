package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type MysqlConf struct {
	DriverName 		string	`yaml:"driverName"`
	DataSourceName 	string	`yaml:"dataSourceName"`
}

func Dbconf(path string)(*MysqlConf,error)  {
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

func Dbinit(path string)(d *sql.DB,err error)  {
	conf,err := Dbconf(path)
	if err!=nil {
		return
	}
	d, err = sql.Open(conf.DriverName,conf.DataSourceName)
	if err!=nil{
		return
	}
	return
}
