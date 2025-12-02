package models

type Posts struct {
	ID        int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id" db:"id"`
	Title     string `gorm:"column:title;not null" json:"title" db:"title"`
	Content   string `gorm:"column:content;not null" json:"content" db:"content"`
	UserID    int64  `gorm:"column:user_id;not null" json:"user_id" db:"user_id"`
	CreatedAt string `gorm:"column:created_at" json:"created_at" db:"created_at"`
	UpdatedAt string `gorm:"column:updated_at" json:"updated_at" db:"updated_at"`
	//gorm.Model
}

func (Posts) TableName() string {
	return "posts"
}
