package main

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id        uint      `gorm:"primaryKey" `
	Username  string    `gorm:"size:50;not null;uniqueIndex"`
	Email     string    `gorm:"size:255;not null;uniqueIndex"`
	Password  string    `gorm:"size:255;not null"`
	Nickname  string    `gorm:"size:50;not null"`
	Avatar    string    `gorm:"size:500;not null"`
	Bio       string    `gorm:"size:500;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt time.Time `gorm:"not null"`

	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Post struct {
	Id        uint      `gorm:"primaryKey" `
	Title     string    `gorm:"size:200;not null"`
	Content   string    `gorm:"type:longtext;not null"`
	Slug      string    `gorm:"size:255;not null"`
	Status    string    `gorm:"size:20;default:draft"` //draft/published/archived
	ViewCount int       `gorm:"default:0"`
	UserId    uint      `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserId"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt time.Time `gorm:"not null"`

	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	Id        uint      `gorm:"primaryKey" `
	Content   string    `gorm:"type:longtext;not null"`
	Author    string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:255;not null"`
	Website   string    `gorm:"size:500;not null"`
	PostId    uint      `gorm:"not null;index"`
	Post      Post      `gorm:"foreignKey:PostId"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt time.Time `gorm:"index"`
}

func main() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := user + ":" + pass + "@tcp(localhost:3306)/" + dbname + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("mysql open err:", err)
		return
	}

	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		fmt.Println("mysql auto migrate err:", err)
		return
	}

}
