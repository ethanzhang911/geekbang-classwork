package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type user struct {
	id   int
	name string
	age  int
}

func Query() (err error) {
	var db *sql.DB
	var u user
	//初始化数据库连接
	dsn := "root:Lucifer123@tcp(127.0.0.1:3306)/user"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//查询
	sqlStr := "select id, name, age from user where id=?"
	err = db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)

	if err == sql.ErrNoRows{
		return fmt.Errorf("%s,%w","空行",err)
	}else {
		return err
	}



}



