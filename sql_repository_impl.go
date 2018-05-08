package ezorm

import (
	"fmt"
	"text/template"
	"errors"
)

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/5/8
 * 
*/

type sqlRepositoryImpl struct {
	tpls map[string]sqlStatement
}

func(f *sqlRepositoryImpl)build(sqlType,id string,param interface{})(sql string,args []interface{},err error){
	if tpl,exist:= f.tpls[id];exist{
		return tpl.prepare(sqlType,param)
	}

	err = fmt.Errorf("CANNOT FIND STATEMENT WITH ID:%v",id)

	return
}

func(f *sqlRepositoryImpl)addStatement(id,sqlType,driverName string,tpl *template.Template)error{
	if id==""{
		return errors.New("STATEMENT ID CANNOT BE EMPTY")
	}

	if _,exist := f.tpls[id];exist{
		return fmt.Errorf("DUPLICATE STATEMENT ID:%v",id)
	}

	f.tpls[id] = newSqlStatement(id,sqlType,driverName,tpl)

	return nil
}
