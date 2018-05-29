package ezorm

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"text/template"
	"bytes"
)

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/5/8
 *
 */

type sqlTemplateImpl struct {
	id         string
	sqlType    string
	driverName string
	tpl        *template.Template
}

func (s *sqlTemplateImpl) prepare(sqlType string,param interface{}) (sql string, args []interface{}, err error) {
	bf := bytes.Buffer{}
	err = s.tpl.Execute(&bf, param)

	if err != nil {
		return
	}

	rawSqlText := string(bf.Bytes())

	if param==nil{
		sql = rawSqlText
		return
	}

	if sqlType!=s.sqlType{
		err = errors.New("SQL TYPE MISMATCH")
		return
	}

	paramPtr := reflect.ValueOf(param)

	if paramPtr.Kind() != reflect.Ptr {
		err = errors.New("STATEMENT ONLY SUPPORT POINT PARAM")
		return
	}

	paramVal := reflect.Indirect(paramPtr)

	if paramVal.Kind() != reflect.Struct {
		err = errors.New("paramVal ONLY SUPPORT STRUCT TYPE")
		return
	}

	r, err := regexp.Compile(`#{[^}]+}`)

	if err != nil {
		return
	}

	argNames := r.FindAllString(rawSqlText, -1)
	args = make([]interface{}, 0, len(argNames))

	for _, argName := range argNames {
		fVal := paramVal.FieldByName(argName[2 : len(argName)-1])

		if fVal.Kind() == reflect.Invalid {
			err = fmt.Errorf("INVALID PARAMETER NAME IN SQL:%v", argName)
			return
		}

		args = append(args, fVal.Interface())
	}

	switch s.driverName {
	case "postgres":
		i := 0
		sql = r.ReplaceAllStringFunc(rawSqlText, func(str string) string {
			i++
			return fmt.Sprintf("$%v", i)
		})
	default:
		sql = r.ReplaceAllString(rawSqlText, "?")
	}

	return
}
