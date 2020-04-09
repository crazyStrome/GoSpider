package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var conn *sql.DB

func init() {
	InitMySQLDB()
}

// InitMySQLDB 实现mysql的实例化处理,单例模式
func InitMySQLDB() {
	if conn == nil {
		log.Println("MySQLDB initializing...")
		conn, _ = sql.Open("mysql", "root:000000@tcp(localhost:3306)/spider_store")
		log.Println("MySQLDB initialized")
	}
}

// Save 把链接和标题储存在mysql数据库中
func Save(url string, title string) (id int64, err error) {
	// InitMySQLDB()
	stmt, err := conn.Prepare("insert into weblinks(url, title) values(?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	rs, err := stmt.Exec(url, title)
	if err != nil {
		log.Println(err)
		return
	}
	id, err = rs.LastInsertId()
	return
}

// Close 用来关闭mysql连接
func Close() {
	if conn != nil {
		conn.Close()
	}
}

// Contet 中保存了MySQL中的表信息
type Content struct {
	Id    int64
	Url   string
	Title string
}

// GetQueryResult returns the cols and rows data
func GetQueryResult() (cols []string, cs []Content) {
	rows, err := conn.Query("select * from weblinks")
	if err != nil {
		log.Println(err)
	}

	cols, err = rows.Columns()
	if err != nil {
		log.Println(err)
	}

	var id int64
	var url string
	var title string

	for rows.Next() {
		rows.Scan(&id, &url, &title)
		c := Content{
			Id:    id,
			Url:   url,
			Title: title,
		}
		cs = append(cs, c)
	}
	return
}
