package ezorm

import "text/template"

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/5/8
 * 
*/

type sqlRepository interface {
	build(sqlType,id string,param interface{})(sql string,args []interface{},err error)
	addStatement(id,sqlType,driverName string,tpl *template.Template)error
}

func newSqlRepository()sqlRepository{
	return &sqlRepositoryImpl{
		tpls:make(map[string]sqlStatement),
	}
}

