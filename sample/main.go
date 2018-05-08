package main

import (
	"github.com/schbook/ezorm"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/5/8
 * 
*/

type RowType struct{
	Id int
	Name string
}

type DAO struct{
	Id int
	Name string
	List []RowType `stmt:"complexQueryRowTest"`
}

func main(){
	err := ezorm.Register("default","mysql","root:@tcp(localhost:3306)/mysql?charset=utf8&parseTime=True&loc=Local","./config/")
	if err!=nil{
		log.Fatal(err.Error())
	}

	rs,err := ezorm.Use("default").Update("main","createTestDB",nil)

	if err!=nil{
		log.Fatal(err.Error())
	}

	log.Printf("createTestDB result:%v\n",rs)

	rs,err = ezorm.Use("default").Update("main","createTestTable",nil)

	if err!=nil{
		log.Fatal(err.Error())
	}

	log.Printf("createTestTable result:%v\n",rs)

	rs,err = ezorm.Use("default").Insert("main","insertTest", &struct {
		Name string
	}{
		Name:"hello ezorm",
	})

	if err!=nil{
		log.Fatal(err.Error())
	}

	log.Printf("insertTest result:%v\n",rs)

	rows := make([]RowType,0)

	err = ezorm.Use("default").QueryRows("main","queryRowsTest",nil,&rows)

	if err!=nil{
		log.Fatal(err.Error())
	}

	log.Printf("queryRowsTest result:%v\n",rows)

	row := RowType{
		Id:1,
	}

	err = ezorm.Use("default").QueryRow("main","queryRowTest",&row,&row)

	if err!=nil{
		log.Fatal(err.Error())
	}

	log.Printf("queryRowTest result:%v\n",row)

	dao := DAO{
		Id:1,
		Name:"hello ezorm",
	}

	err = ezorm.Use("default").Query(&dao)

	if err!=nil{
		log.Fatal(err.Error())
	}

	log.Printf("complexQueryRowTest result:%v\n",dao)
}
