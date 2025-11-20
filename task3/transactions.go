package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Accounts struct {
	id      int    `gorm:"column:id;primaryKey" json:"id"`
	name    string `gorm:"column:name;not null" json:"name"`
	balance int    `gorm:"column:balance;not null" json:"balance"`
}

type Transactions struct {
	id              int `gorm:"column:id;primaryKey" json:"id"`
	from_account_id int `gorm:"column:from_account_id;not null" json:"from_account_id"`
	to_account_id   int `gorm:"column:to_account_id;not null" json:"to_account_id"`
	amount          int `gorm:"column:amount;not null" json:"amount"`
}

func connectDB() (*gorm.DB, error) {
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

	return db, err
}

func addAccount(ctx context.Context, mysqlDB *sql.DB, name string, balance int) error {
	//1.开启事务
	tx, err := mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("tx begin err: %v", err)
	}
	//确保函数退出时，没有得及回滚的操作
	defer func() {
		p := recover()
		if p != nil {
			tx.Rollback()
		}
	}()

	//2.新增账户信息
	_, err = tx.ExecContext(ctx, "INSERT INTO accounts(name, balance) VALUES(?,?)", name, balance)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("tx err: %v", err)
	}
	//3.提交事务
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx commit err: %v", err)
	}
	return nil
}

func getAccountId(ctx context.Context, mysqlDB *sql.DB, name string) (int, error) {
	//1.开启事务
	tx, err := mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("tx begin err: %v", err)
	}
	//确保函数退出时，没有得及回滚的操作
	defer func() {
		p := recover()
		if p != nil {
			tx.Rollback()
		}
	}()

	//2.查询数据库获取数据
	findId := 0
	err = tx.QueryRowContext(ctx, "SELECT id FROM accounts WHERE name = ?", name).Scan(&findId)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("tx err: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("tx commit err: %v", err)
	}
	return findId, nil
	//3.提交事务
}

func transfer(ctx context.Context, mysqlDB *sql.DB, from_id, to_id int, amount int) error {

	//1.开启事务
	tx, err := mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Transfer BeginTx err:%w", err)
	}

	//确保函数退出时，没有得及回滚的操作
	defer func() {
		p := recover()
		if p != nil {
			tx.Rollback()
		}
	}()

	//2.查询账户，确保余额充足
	fromBalance := 0
	err = tx.QueryRowContext(ctx, "select balance from accounts where id=?", from_id).Scan(&fromBalance)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Transfer QueryRowContext balance err:%w", err)
	}

	if fromBalance < amount {
		tx.Rollback()
		return fmt.Errorf("Transfer from balance %d not enough", fromBalance)
	}

	//3. 转账：扣除当前账户余额
	_, err = tx.ExecContext(ctx, "update accounts set balance=? where id=?", fromBalance-amount, from_id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Transfer Sub balance err:%w", err)
	}

	//4. 入账，更新转入账户余额
	oldBalance := 0
	err = tx.QueryRowContext(ctx, "select balance from accounts where id=?", to_id).Scan(&oldBalance)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Transfer get target accounter old balance err:%w", err)
	}

	_, err = tx.ExecContext(ctx, "update accounts set balance=? where id=?", oldBalance+amount, to_id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Transfer add amount to old balance err:%w", err)
	}

	//5. 将转账流水登记下来
	_, err = tx.ExecContext(ctx, "insert into transactions (from_account_id,to_account_id,amount) values (?,?,?)",
		from_id, to_id, amount)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Transfer record  err:%w", err)
	}

	//6.提交事务
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Transfer commit err:%w", err)
	}

	return nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		fmt.Println("connectDB err:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("获取底层 DB 失败: %v", err)
		return
	}
	if err := sqlDB.Ping(); err != nil {
		fmt.Println("数据库 Ping 失败: %v", err)
		return
	}

	defer sqlDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addAccount(ctx, sqlDB, "A", 90)
	addAccount(ctx, sqlDB, "B", 50)

	from_id, err := getAccountId(ctx, sqlDB, "A")
	if err != nil {
		fmt.Println("<UNK> getAccount fromId err:", err)
	}

	to_id, err := getAccountId(ctx, sqlDB, "B")
	if err != nil {
		fmt.Println("<UNK> getAccount toId err:", err)
	}

	err = transfer(ctx, sqlDB, from_id, to_id, 100)
	if err != nil {
		fmt.Println("<UNK> transfer err:", err)
	}

}
