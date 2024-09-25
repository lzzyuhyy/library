package consts

const (
	MySQL  = "mysql"
	ES     = "es"
	AliPay = "alipay"
)

var IndexName = []string{"books", "users", "records"}

var MysqlConf struct {
	Host   string
	Port   int64
	User   string
	Pass   string
	Dbname string
}

var EsConf struct {
	Address string
}

var AlipayConf struct {
	Appid      string `json:"Appid"`
	PrivateKey string `json:"PrivateKey"`
	PublicKey  string `json:"PublicKey"`
	Subject    string `json:"Subject"`
	ReturnURL  string `json:"ReturnURL"`
	NotifyURL  string `json:"NotifyURL"`
}
