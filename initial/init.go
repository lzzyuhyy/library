package initial

import (
	"encoding/json"
	"gorm.io/gorm"
	"library/consts"
	"library/consul"
	"library/db"
	"library/models/books"
	"library/models/records"
	"library/models/user"
	"library/pkg/es"
	"log"
	"strconv"
)

func initReadConsulConf() error {
	conf, err := consul.GetKV(consts.MySQL)
	if err != nil {
		return err
	}

	err = json.Unmarshal(conf, &consts.MysqlConf)
	if err != nil {
		return err
	}

	esconf, err := consul.GetKV(consts.ES)
	if err != nil {
		return err
	}

	err = json.Unmarshal(esconf, &consts.EsConf)
	if err != nil {
		return err
	}

	payconf, err := consul.GetKV(consts.AliPay)
	if err != nil {
		return err
	}

	err = json.Unmarshal(payconf, &consts.AlipayConf)
	if err != nil {
		return err
	}

	return nil
}

func autoMigrate() error {
	return db.WithMysql(func(db *gorm.DB) error {
		return db.AutoMigrate(
			new(books.Book),
			new(user.User),
			new(records.Record),
		)
	})
}

func initEsIndex() error {
	return es.CreateIndex(consts.IndexName...)
}

func InitEsData() error {
	getBooks, err := books.GetBooks()
	if err != nil {
		return err
	}
	for _, v := range getBooks {
		marshal, err := json.Marshal(&v)
		if err != nil {
			log.Println("同步es，转码失败：", err)
			continue
		}

		err = es.SyncData("books", strconv.Itoa(int(v.ID)), marshal)
		if err != nil {
			log.Println("同步es失败：", err)
			continue
		}
	}

	return nil
}

func Initial() error {
	err := initReadConsulConf()
	if err != nil {
		return err
	}

	err = autoMigrate()
	if err != nil {
		return err
	}

	err = initEsIndex()
	if err != nil {
		return err
	}

	go InitEsData()

	return nil
}
