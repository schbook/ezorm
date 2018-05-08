# EZORM

An easy ORM tool for Golang, support MyBatis-Like XML template SQL

## Overview

* Full-Featured ORM (almost)
* MyBatis-Like XML template SQL
* Developer Friendly

## Getting Started

```go
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
```

## License

© Schbook, 2018~time.Now

Released under the [MIT License](https://github.com/schbook/ezorm/blob/master/License)