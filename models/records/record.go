package records

import (
	"gorm.io/gorm"
	"library/db"
	"time"
)

type Record struct {
	gorm.Model
	Uid     uint       `gorm:"type:int(11);not null;comment:用户id"`
	BookId  uint       `gorm:"type:int(11);not null;comment:图书id"`
	EndDate *time.Time `gorm:"type:datetime;not null;comment:归还日期"`
}

func Add(r *Record) error {
	return db.WithMysql(func(db *gorm.DB) error {
		return db.Create(r).Error
	})
}

func All() (list []Record, err error) {
	err = db.WithMysql(func(db *gorm.DB) error {
		return db.Find(&list).Error
	})
	return
}

func UpdateEndDate(id uint, newDate *time.Time) error {
	return db.WithMysql(func(db *gorm.DB) error {
		return db.Model(&Record{}).Where("id = ?", id).Update("end_date", newDate).Error
	})
}
