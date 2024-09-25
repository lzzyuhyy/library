package books

import (
	"gorm.io/gorm"
	"library/db"
	"time"
)

type Book struct {
	gorm.Model
	Title      string     `gorm:"type:varchar(255);not null;comment:名"`
	Image      string     `gorm:"type:varchar(255);not null;comment:图"`
	Type       int8       `gorm:"type:tinyint(1);not null;comment:类别"`
	Author     string     `gorm:"type:varchar(255);not null;comment:作者"`
	Isbn       string     `gorm:"type:varchar(255);not null;comment:isbn编码"`
	SaleDate   *time.Time `gorm:"type:datetime;not null;comment:出版日期"`
	Popularity int64      `gorm:"type:bigint;not null;comment:名"`
}

func GetBooks() (list []Book, err error) {
	err = db.WithMysql(func(db *gorm.DB) error {
		return db.Find(&list).Error
	})
	return
}

func GetBookById(id uint) (book Book, err error) {
	err = db.WithMysql(func(db *gorm.DB) error {
		return db.Limit(1).Find(&book, id).Error
	})
	return
}
