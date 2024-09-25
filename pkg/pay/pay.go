package pay

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"library/consts"
	"net/http"
)

func newAlipayClient() (*alipay.Client, error) {
	return alipay.NewClient(consts.AlipayConf.Appid, consts.AlipayConf.PrivateKey, false)
}

func TradePagePay(ctx context.Context, tradeNo, totalAmount, subject string) (string, error) {
	cli, err := newAlipayClient()
	if err != nil {
		return "", err
	}

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", tradeNo).
		Set("total_amount", totalAmount).
		Set("subject", subject).
		Set("return_url", consts.AlipayConf.ReturnURL).
		Set("notify_url", consts.AlipayConf.NotifyURL)

	return cli.TradePagePay(ctx, bm)
}

func VerifySign(r *http.Request) error {
	bodyMap, err := alipay.ParseNotifyToBodyMap(r)
	if err != nil {
		return err
	}

	sign, err := alipay.VerifySign(consts.AlipayConf.PublicKey, bodyMap)
	if err != nil || !sign {
		return fmt.Errorf("验签失败")
	}

	return nil
}
