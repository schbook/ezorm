package ezorm

import (
	"encoding/xml"
	"io/ioutil"
	"text/template"
	"strings"
	"errors"
	"reflect"
	"database/sql"
	"fmt"
)

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/5/8
 * 
*/

type sqlFactoryImpl struct {
	tplRepositories map[string]sqlRepository
	dataSourceName string
	driverName string
	dbName string
}

type mapperXML struct{
	Namespace     string `xml:"namespace,attr"`
	Children []templateElem `xml:",any"`
}

type templateElem struct{
	XMLName xml.Name
	Id string `xml:"id,attr""`
	Sql string `xml:",chardata"`
}

func(fac *sqlFactoryImpl) build(dbName,driverName,dataSourceName,mapperPath string) error{
	fac.tplRepositories = make(map[string]sqlRepository)
	fac.dbName = dbName
	fac.dataSourceName = dataSourceName
	fac.driverName = driverName

	files, err := ioutil.ReadDir(mapperPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir(){
			continue
		}

		filePath := mapperPath+"/"+f.Name()

		err = fac.parseMapperFile(driverName,filePath)

		if err!=nil{
			return err
		}
	}

	return nil
}

func(fac *sqlFactoryImpl) parseMapperFile(driverName,mapperFilePath string)error{
	content, err := ioutil.ReadFile(mapperFilePath)
	if err != nil {
		return err
	}

	var mapper mapperXML
	err = xml.Unmarshal(content, &mapper)
	if err != nil {
		return err
	}

	rep,exist := fac.tplRepositories[mapper.Namespace]

	if !exist{
		rep = newSqlRepository()
		fac.tplRepositories[mapper.Namespace] = rep
	}

	for _,m:=range mapper.Children{
		tpl,err := template.New(m.Id).Parse(m.Sql)
		if err!=nil{
			return err
		}

		err = rep.addStatement(m.Id,strings.ToLower(m.XMLName.Local),driverName,tpl)

		if err!=nil{
			return err
		}
	}

	return nil
}

func(fac *sqlFactoryImpl) QueryRows(namespace,id string,param interface{},rows interface{})error{
	if fac.tplRepositories == nil{
		return fmt.Errorf("SQL FACTORY %v NOT FOUND",fac.dbName)
	}

	rep,exist := fac.tplRepositories[namespace]

	if !exist{
		return errors.New("NAMESPACE NOT FOUND")
	}

	rowsPtr := reflect.ValueOf(rows)

	if rowsPtr.Kind()!=reflect.Ptr{
		return errors.New("ROWS ONLY SUPPORT POINT TYPE")
	}

	rowsVal := reflect.Indirect(rowsPtr)

	if rowsVal.Kind()!=reflect.Slice{
		return errors.New("THE VALUE OF ROWS ONLY SUPPORT SLICE TYPE")
	}

	sql,args,err := rep.build("select",id,param)

	if err!=nil{
		return err
	}

	db,err := fac.GetDB()

	if err!=nil{
		return err
	}

	rs,err := db.Query(sql,args...)

	if err!=nil{
		return err
	}

	defer rs.Close()

	colsMap,scanList,err := getColsInfo(rs)

	if err!=nil{
		return err
	}

	rowType := rowsVal.Type().Elem()

	for rs.Next(){
		item := reflect.New(rowType).Elem()

		err = rs.Scan(scanList...)

		if err!=nil{
			return err
		}

		setRowVal(rowType,item,colsMap,scanList)

		rowsVal.Set(reflect.Append(rowsVal,item))
	}

	return nil
}

func(fac *sqlFactoryImpl) QueryRow(namespace,id string,param interface{},row interface{})error{
	if fac.tplRepositories == nil{
		return fmt.Errorf("SQL FACTORY %v NOT FOUND",fac.dbName)
	}

	rep,exist := fac.tplRepositories[namespace]

	if !exist{
		return errors.New("NAMESPACE NOT FOUND")
	}

	rowsPtr := reflect.ValueOf(row)

	if rowsPtr.Kind()!=reflect.Ptr{
		return errors.New("ROW ONLY SUPPORT POINT TYPE")
	}

	sql,args,err := rep.build("select",id,param)

	if err!=nil{
		return err
	}

	db,err := fac.GetDB()

	if err!=nil{
		return err
	}

	rs,err := db.Query(sql,args...)

	if err!=nil{
		return err
	}

	defer rs.Close()

	colsMap,scanList,err := getColsInfo(rs)

	if err!=nil{
		return err
	}

	rowVal := reflect.Indirect(rowsPtr)
	rowType := rowVal.Type()

	if rs.Next(){
		item := reflect.New(rowType).Elem()

		err = rs.Scan(scanList...)

		if err!=nil{
			return err
		}

		setRowVal(rowType,item,colsMap,scanList)

		rowVal.Set(item)
	}

	return nil
}

func(fac *sqlFactoryImpl)Query(param interface{}) error{
	if fac.tplRepositories == nil{
		return fmt.Errorf("SQL FACTORY %v NOT FOUND",fac.dbName)
	}

	paramPtr := reflect.ValueOf(param)

	if paramPtr.Kind()!=reflect.Ptr{
		return errors.New("QUERY ONLY SUPPORT POINT PARAM")
	}

	paramVal := reflect.Indirect(paramPtr)
	paramType := paramVal.Type()
	ns := paramType.PkgPath()

	for i:=0;i<paramType.NumField();i++{
		f := paramType.Field(i)
		stmtId := f.Tag.Get("stmt")

		if stmtId==""{
			continue
		}

		var err error

		fVal := paramVal.Field(i).Addr().Interface()

		if f.Type.Kind() == reflect.Slice {
			err = fac.QueryRows(ns,stmtId,param,fVal)
		}else {
			err = fac.QueryRow(ns,stmtId,param,fVal)
		}

		if err!=nil{
			return err
		}
	}

	return nil
}

func(fac *sqlFactoryImpl) Update(namespace string, id string, param interface{})(rowsAffected int64,err error){
	return fac.Execute("update",namespace,id,param)
}

func(fac *sqlFactoryImpl) Insert(namespace string, id string, param interface{})(lastInsertId int64,err error){
	return fac.Execute("insert",namespace,id,param)
}

func(fac *sqlFactoryImpl) Execute(sqlType, namespace string, id string, param interface{})(cnt int64,err error){
	if fac.tplRepositories == nil{
		err = fmt.Errorf("SQL FACTORY %v NOT FOUND",fac.dbName)
		return
	}

	rep,exist := fac.tplRepositories[namespace]

	if !exist{
		err = errors.New("NAMESPACE NOT FOUND")

		return
	}

	sql,args,err := rep.build(sqlType,id,param)

	if err!=nil{
		return
	}

	db,err := fac.GetDB()

	if err!=nil{
		return
	}

	defer db.Close()

	rs,err := db.Exec(sql,args...)

	if err!=nil{
		return
	}

	if sqlType=="insert"{
		cnt,err = rs.LastInsertId()
	}else{
		cnt,err = rs.RowsAffected()
	}

	return
}

func(fac *sqlFactoryImpl)GetDB()(db *sql.DB, err error){
	if fac.driverName==""||fac.dataSourceName==""{
		return nil,fmt.Errorf("SQL FACTORY %v NOT FOUND",fac.dbName)
	}

	return sql.Open(fac.driverName, fac.dataSourceName)
}