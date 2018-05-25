package ezorm

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/5/8
 * 
*/

var(
	factories = make(map[string]SqlFactory)
)

func Register(dbName,driverName,dataSourceName,mapperPath string)error{
	f,err := newSqlFactory(dbName, driverName,dataSourceName,mapperPath)

	if err!=nil{
		return err
	}

	factories[dbName] = f

	return nil
}

func Use(dbName string)SqlFactory{
	return factories[dbName]
}