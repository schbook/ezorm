package ezorm

import "text/template"

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/5/8
 * 
*/

type sqlStatement interface {
	prepare(sqlType string,param interface{})(sql string,args []interface{},err error)
}

func newSqlStatement(id,sqlType,driverName string,tpl *template.Template) sqlStatement{
	return &sqlTemplateImpl{
		id:id,
		sqlType:sqlType,
		driverName:driverName,
		tpl:tpl,
	}
}

