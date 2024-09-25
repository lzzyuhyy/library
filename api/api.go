package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"library/models/books"
	"library/models/records"
	"library/pkg/es"
	"library/pkg/pay"
	"log"
	"net/http"
	"time"
)

type bookListQuery struct {
	Query struct {
		Match struct {
			Title string `json:"Title"`
		} `json:"match"`
	} `json:"query"`
	From int `json:"from"`
	Size int `json:"size"`
	Sort []struct {
		SaleDate struct {
			Order string `json:"order"`
		} `json:"SaleDate"`
	} `json:"sort"`
}

type bookListRes struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{} `json:"max_score"`
		Hits     []struct {
			Index  string      `json:"_index"`
			Type   string      `json:"_type"`
			Id     string      `json:"_id"`
			Score  interface{} `json:"_score"`
			Source struct {
				ID         int         `json:"ID"`
				CreatedAt  time.Time   `json:"CreatedAt"`
				UpdatedAt  time.Time   `json:"UpdatedAt"`
				DeletedAt  interface{} `json:"DeletedAt"`
				Title      string      `json:"Title"`
				Image      string      `json:"Image"`
				Type       int         `json:"Type"`
				Author     string      `json:"Author"`
				Isbn       string      `json:"Isbn"`
				SaleDate   time.Time   `json:"SaleDate"`
				Popularity int         `json:"Popularity"`
			} `json:"_source"`
			Sort []int64 `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

type Source struct {
	ID         int         `json:"id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  interface{} `json:"deleted_at"`
	Title      string      `json:"title"`
	Image      string      `json:"image"`
	Type       int         `json:"type"`
	Author     string      `json:"author"`
	Isbn       string      `json:"isbn"`
	SaleDate   time.Time   `json:"sale_date"`
	Popularity int         `json:"popularity"`
}

func BookList(c *gin.Context) {
	var req bookListQuery
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "params error",
		})
		return
	}

	// logic
	marshal, err := json.Marshal(&req)
	if err != nil {
		log.Println("json encoding error", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "json encoding error",
		})
		return
	}
	data, err := es.SearchData("books", string(marshal))
	if err != nil {
		log.Println("get data error", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get data error",
		})
		return
	}
	var res bookListRes
	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Println("json encoding error", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "json encoding error",
		})
		return
	}

	var info []Source
	for _, v := range res.Hits.Hits {
		var item Source
		err = copier.Copy(&item, &v.Source)
		if err != nil {
			log.Println("struct copy error", err)
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "struct copy error",
			})
			return
		}

		info = append(info, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": info,
	})
	return
}

type borrowBookReq struct {
	Uid         uint  `json:"uid"`
	BookId      uint  `json:"book_id"`
	EndDuration int64 `json:"end_duration"`
}

func BorrowBook(c *gin.Context) {
	var req borrowBookReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "params error",
		})
		return
	}
	endTime := time.Now().Add(time.Duration(req.EndDuration*24) * time.Hour)
	err := records.Add(&records.Record{
		Uid:     req.Uid,
		BookId:  req.BookId,
		EndDate: &endTime,
	})
	if err != nil {
		log.Println("add record error", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "add record error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
	return
}

type payReq struct {
	BookId      uint  `json:"book_id"`
	Num         int64 `json:"num"`
	EndDuration int64 `json:"end_duration"` // 3天一元--不足按三天算
}

func Pay(c *gin.Context) {
	var req borrowBookReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "params error",
		})
		return
	}

	info, err := books.GetBookById(req.BookId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "data not find",
		})
		return
	}

	amount := decimal.NewFromFloat(10 / 3.0).Ceil().String()
	no := uuid.NewString()
	payUrl, err := pay.TradePagePay(c, no, amount, info.Title)
	if err != nil {
		log.Println("pay url generate error", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "pay url generate error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": payUrl,
	})
	return
}

func Notify(c *gin.Context) {
	err := pay.VerifySign(c.Request)
	if err != nil {
		log.Println(err)
		return
	}

	c.Writer.Write([]byte("success"))
}
