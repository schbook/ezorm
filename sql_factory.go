package ezorm

import "database/sql"

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/5/8
 * 
*/

type SqlFactory interface {
	QueryRows(namespace,id string,param interface{},rows interface{})error
	QueryRow(namespace,id string,param interface{},row interface{})error
	Query(param interface{}) error
	Insert(namespace string, id string, param interface{})(lastInsertId int64,err error)
	Update(namespace string, id string, param interface{})(rowsAffected int64,err error)
	Execute(sqlType, namespace string, id string, param interface{})(cnt int64,err error)
	GetDB()(db *sql.DB, err error)
}

func newSqlFactory(dbName,driverName,dataSourceName,mapperPath string) (SqlFactory,error){
	f := sqlFactoryImpl{}
	err := f.build(dbName,driverName,dataSourceName,mapperPath)

	if err!=nil{
		return nil,err
	}

	return &f,nil
}