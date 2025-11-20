package main

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employee struct {
	Id         int     `gorm:"column:id;primaryKey" json:"id"`
	Name       string  `gorm:"column:name;not null" json:"name"`
	Department string  `gorm:"column:department;not null" json:"department"`
	Salary     float32 `gorm:"column:salary;not null" json:"salary"`
}

func main() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	users := []Employee{
		{1, "张三", "技术部", 9600},
		{2, "李四", "市场部", 8000},
		{3, "王五", "技术部", 15000},
		{4, "赵,,六", "宣传部", 6000},
	}

	dsn := user + ":" + pass + "@tcp(127.0.0.1:3306)/" + dbname + "?charset=utf8mb4&parseTime=true"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("sqlx connect error: ", err)
		return
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//1.插入数据
	for _, u := range users {
		_, err = db.ExecContext(ctx, "insert into employees(name,department,salary) values(?,?,?)",
			u.Name, u.Department, u.Salary)
		if err != nil {
			fmt.Println("insert error: ", err)
			return
		}
	}

	//2.查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	var emps []Employee
	sqlcmd := `select * from employees where department=?`

	const dept = "技术部"
	err = db.SelectContext(ctx, &emps, sqlcmd, dept)
	if err != nil {
		fmt.Println("select error: ", err)
	}

	for _, emp := range emps {
		fmt.Printf("\rid:%d,name:%s department:%s,salary:%2.f\n ",
			emp.Id, emp.Name, emp.Department, emp.Salary)
	}

	//3.使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。,
	var emp Employee
	sqlcmd = `select * from employees order by salary desc limit 1`
	err = sqlx.GetContext(ctx, db, &emp, sqlcmd)
	if err != nil {
		fmt.Println("select error: ", err)
	} else {
		fmt.Printf("\rid:%d,name:%s department:%s,salary:%2.f\n ",
			emp.Id, emp.Name, emp.Department, emp.Salary)
	}

}
