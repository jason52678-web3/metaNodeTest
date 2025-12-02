package models

type Comments struct {
	ID        int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id" db:"id"`
	Content   string `gorm:"column:content;not null" json:"content" db:"content"`
	UserID    int64  `gorm:"column:user_id;not null" json:"user_id" db:"user_id"`
	PostID    int    `gorm:"column:post_id;not null" json:"post_id" db:"post_id"`
	CreatedAt string `gorm:"column:created_at" json:"created_at" db:"created_at"`
}

func (Comments) TableName() string {
	return "comments"
}
