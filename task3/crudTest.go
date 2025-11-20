package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || name == "" {
		log.Fatal("请设置环境变量：DB_USER、DB_PASSWORD、DB_NAME")
	}
	fmt.Println(user, pass, name)

	// 注意：parseTime=True&loc=Local 用于正确解析 MySQL 的 DATETIME/TIMESTAMP
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, name)
	fmt.Println("dsn: ", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
		return nil
	}

	return db
}

type Student struct {
	name   string `gorm:"column:name;not null" json:"name"`
	gender bool   `gorm:"column:gender" json:"gender"`
	age    int    `gorm:"column:age" json:"age"`
	grade  string `gorm:"column:grade" json:"grade"`
	//gorm.Model
}

func main() {
	db := connectDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层 DB 失败: %v", err)
		return
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("数据库 Ping 失败: %v", err)
		return
	}
	log.Println("✅ 数据库连接成功")

	//1.编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级",性别为男。,
	_, err = sqlDB.Exec("INSERT INTO students (name,gender,age,grade) VALUES (?,?,?,?)",
		"张三", true, 20, "三年级")

	if err != nil {
		log.Fatalf("insert failed: %v", err) // ✅ 正确打印错误信息
	} else {
		log.Println("insert success")
	}

	//2.编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。,
	rows, err := sqlDB.Query("select * from students where age>?", 18)
	if err != nil {
		log.Fatalf("query failed: %v", err) // ✅ 正确打印错误信息
	} else {
		log.Println("query success")
	}
	//显示列(字段名)信息
	cols, err := rows.Columns()
	if err != nil {
		log.Fatalf("get columns failed: %v", err)
	} else {
		log.Printf("columns: %v", cols)
	}

	//显示行(记录)信息
	for rows.Next() {
		var student Student
		id := 0
		err := rows.Scan(&id, &student.name, &student.gender, &student.age, &student.grade)
		if err != nil {
			log.Fatalf("scan failed: %v", err)
		}
		log.Printf("student:%d , %v", id, student)
	}

	//3.编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。,
	_, err = sqlDB.Exec("update students set grade=? where name=?",
		"四年级", "张三")
	if err != nil {
		log.Fatalf("update failed: %v", err)
	} else {
		log.Println("update success")
	}

	//4.编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	_, err = sqlDB.Exec("DELETE FROM students WHERE age<?", 15)
	if err != nil {
		log.Fatalf("delete failed: %v", err)
	} else {
		log.Println("delete success")
	}

}
