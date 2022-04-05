package database

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// Route User 路由表结构体
type Route struct {
	Id             string `db:"id"`
	Host           string `db:"host"`
	Port           int64  `db:"port"`
	Path           string `db:"path"`
	HttpMethod     string `db:"http_method"`
	HttpTemplateId string `db:"http_template_id"`
}

// QueryAllRoutes 查询所有路由信息
func QueryAllRoutes() ([]Route, error) {
	routes := make([]Route, 0)
	rows, err := MysqlDb.Query("select id, host, port, path, http_method, http_template_id from route")
	fmt.Println(err)
	var route Route
	for rows.Next() {
		err := rows.Scan(&route.Id, &route.Host, &route.Port, &route.Path, &route.HttpMethod, &route.HttpTemplateId)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func QueryHostRoutes(host string) (int64, error) {
	var count int64
	rows, err := MysqlDb.Query("SELECT * FROM route WHERE `host` =?", host) // Query
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		count += 1
	}
	return count, nil
}

// QueryMockById 根据ID查询单条数据
func QueryMockById(id string) (*Route, error) {
	route := new(Route)
	row := MysqlDb.QueryRow("select * from route where id=?", id)
	if err := row.Scan(&route.Id, &route.Host, &route.Port, &route.Path, &route.HttpMethod, &route.HttpTemplateId); err != nil {
		fmt.Printf("scan failed, err:%v", err)
		return nil, err
	}
	return route, nil
}

// QueryPageLimitByHostOrPath 根据条件查询分页数据
func QueryPageLimitByHostOrPath(host string, path string, limit int64) ([]Route, error) {
	where := "where "
	args := make([]interface{}, 0)
	if host != "" {
		where = where + "host = ? "
		args = append(args, host)
	}
	if path != "" {
		where = where + "path = ?"
		args = append(args, path)
	}
	args = append(args, limit)
	// 通过切片存储
	routes := make([]Route, 0)
	rows, _ := MysqlDb.Query("SELECT * FROM `route`"+where+" limit ?", args)
	// 遍历
	var route Route
	for rows.Next() {
		err := rows.Scan(&route.Id, &route.Host, &route.Port, &route.Path, &route.HttpMethod, &route.HttpTemplateId)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}

// Insert 插入数据 返回：主键id,影响行数,错误信息
func (route *Route) Insert() (string, int64, error) {
	id := uuid.NewV4().String()
	ret, err := MysqlDb.Exec("insert INTO mock(id,host,port,path,http_method,http_template_id) values(?,?,?,?,?,?,?)",
		id, route.Host, route.Port, route.Path, route.HttpMethod, route.HttpTemplateId)

	//影响行数
	rowsAffected, _ := ret.RowsAffected()
	fmt.Println("RowsAffected:", rowsAffected)

	return id, rowsAffected, err
}

// Update 更新数据 返回影响行数
func (route *Route) Update() int64 {
	ret, _ := MysqlDb.Exec("UPDATE mock set host=?,port=?,path=?,http_method=?,http_template_id=? where id=?",
		route.Host, route.Port, route.Path, route.HttpMethod, route.HttpTemplateId, route.Id)

	rowsAffected, _ := ret.RowsAffected()
	fmt.Println("RowsAffected:", rowsAffected)
	return rowsAffected
}

// Delete 删除数据
func Delete(id int64) int64 {

	ret, _ := MysqlDb.Exec("delete from route where id=?", id)
	rowsAffected, _ := ret.RowsAffected()

	fmt.Println("RowsAffected:", rowsAffected)
	return rowsAffected
}

// TableInfo 查询schema下的所有表
func TableInfo(dbName string) map[string]string {
	sqlStr := `SELECT table_name tableName,TABLE_COMMENT tableDesc
			FROM INFORMATION_SCHEMA.TABLES 
			WHERE UPPER(table_type)='BASE TABLE'
			AND LOWER(table_schema) = 'mock'
			ORDER BY table_name asc`

	var result = make(map[string]string)

	rows, err := MysqlDb.Query(sqlStr, dbName)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var tableName, tableDesc string
		err = rows.Scan(&tableName, &tableDesc)
		if err != nil {
			panic(err)
		}
		if len(tableDesc) == 0 {
			tableDesc = tableName
		}
		result[tableName] = tableDesc
	}
	return result
}

// StructTx 事务处理,结合预处理
func StructTx() {

	//事务处理
	tx, _ := MysqlDb.Begin()

	// 新增
	userAddPre, _ := MysqlDb.Prepare("insert into users(name, age) values(?, ?)")
	addRet, _ := userAddPre.Exec("zhaoliu", 15)
	ins_nums, _ := addRet.RowsAffected()

	// 更新
	userUpdatePre1, _ := tx.Exec("update users set name = 'zhansan'  where name=?", "张三")
	upd_nums1, _ := userUpdatePre1.RowsAffected()
	userUpdatePre2, _ := tx.Exec("update users set name = 'lisi'  where name=?", "李四")
	upd_nums2, _ := userUpdatePre2.RowsAffected()

	fmt.Println(ins_nums)
	fmt.Println(upd_nums1)
	fmt.Println(upd_nums2)

	if ins_nums > 0 && upd_nums1 > 0 && upd_nums2 > 0 {
		tx.Commit()
	} else {
		tx.Rollback()
	}

}

// RawQueryField 查询数据，指定字段名,不采用结构体
func RawQueryField() {

	rows, _ := MysqlDb.Query("select id,name from users")
	if rows == nil {
		return
	}
	id := 0
	name := ""
	fmt.Println(rows)
	fmt.Println(rows)
	for rows.Next() {
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
}

// RawQueryAllField 查询数据,取所有字段,不采用结构体
func RawQueryAllField() {

	//查询数据，取所有字段
	rows2, _ := MysqlDb.Query("select * from users")

	//返回所有列
	cols, _ := rows2.Columns()

	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))

	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	result := make(map[int]map[string]string)
	for rows2.Next() {
		//填充数据
		rows2.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	fmt.Println(result)
}
