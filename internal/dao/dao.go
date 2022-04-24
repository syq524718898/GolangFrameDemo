package dao

import (
	"database/sql"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Student struct {
	Id int `json:"id"`
	S_name string `json:"s_name"`
	D_id int `json:"d_id"`
}



type Dao interface {
	Close()
	GetStudent(id int)(*Student,error)
}


type dao struct {
	db *sql.DB
	redis redis.Conn
}

func (d *dao) Close()  {
	d.db.Close()
	d.redis.Close()
}

func Init(dbconf string,redisconf string) (db *sql.DB,con redis.Conn,err error) {
	db, err = Dbinit(dbconf)
	if err!=nil{
		return
	}
	con, err = RedisInit(redisconf)
	if err!=nil{
		return
	}
	return
}

func (d *dao)GetStudent(id int)(*Student,error)  {
	b,err := d.redis.Do("GET","stu_"+strconv.Itoa(id))
	if err != nil {
		return nil,err
	}
	var stu Student
	if b!=nil{
		var stubyte []byte
		stubyte, err = redis.Bytes(b, err)
		if err!=nil{
			return nil,err
		}
		_ = json.Unmarshal(stubyte,&stu)
	}else {
		// cache miss
		rows,err :=d.db.Query("select * from student where id="+strconv.Itoa(id))
		if err!=nil{
			return nil,err
		}
		for rows.Next(){
			rows.Scan(&stu.Id,&stu.S_name,&stu.D_id,)
			//fmt.Println(s)
		}
		rows.Close()
	}
	return &stu,nil
}

func New(db *sql.DB,r redis.Conn)(Dao,error)  {
	var d Dao
	d = &dao{db:db,redis:r}
	return d,nil
}
