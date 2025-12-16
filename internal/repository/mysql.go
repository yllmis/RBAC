package repository

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func NewMysql() {
	fmt.Println("正在连接数据库......")
	my := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"258369",
		"10.148.126.209",
		3306,
		"rbac")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败,请检查参数:", err)
		panic(err)
	}
	Conn = conn
	fmt.Println("数据库连接成功!")

}

func Close() {
	DB, err := Conn.DB()
	if err != nil {
		fmt.Println("数据库关闭失败,请检查:", err)
		return
	}
	DB.Close()
	fmt.Println("数据库连接已关闭!")
}
