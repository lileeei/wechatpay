package wechatpay

import (
	"fmt"
	"testing"
)

const (
	AppId = "wxefe5a7d37bd36609"
	//APP_SECRET = "12580ed29fa76617fb099efcd565802f"
	//APP_STATE  = "h4chb7gG"
	ApiKey   = "ixrAp0o6dY1gek3KcA6Goo16wKPUydwq"
	MchId    = "1480318092"
	CertFile = "/home/lilei/tls/apiclient_cert.pem"
	KeyFile  = "/home/lilei/tls/apiclient_key.pem"
	RootFile = "/home/lilei/tls/rootca.pem"
)

func TestMktTransfer(t *testing.T) {
	openId := "o0Vi71A4stlx3wngX0sosRj-UGCY"

	ps := make(Params)

	tradeId := UUID(18)
	ps.SetString("partner_trade_no", tradeId)
	ps.SetString("openid", openId)
	ps.SetInt("amount", int64(100))
	ps.SetString("desc", "hello, world")
	ps.SetString("spbill_create_ip", "116.211.146.193")

	c := NewClient(AppId, ApiKey, MchId)

	resp, err := c.MktTransfer(ps, CertFile, KeyFile, RootFile)
	if err != nil {
		fmt.Printf("MktTransfer Error %v. \n", err)
		return
	}

	fmt.Printf("resp: %#v. \n", resp)

	fmt.Printf("payment_time = %v. \n", resp["payment_time"])

	return
}

func TestUnifiedOrderJSAPI(t *testing.T) {
	ps := make(Params)

	tradeId := UUID(18)

	ps.SetString("body", "hello, world")
	ps.SetString("out_trade_no", tradeId)
	ps.SetInt("total_fee", int64(100))
	ps.SetString("spbill_create_ip", "116.211.146.193")
	ps.SetString("notify_url", "")
	ps.SetString("openid", "")

	c := NewClient(AppId, ApiKey, MchId)

	resp, err := c.UnifiedOrderJSAPI(ps)
	if err != nil {
		fmt.Printf("UnifiedOrderJSAPI Error %v. \n", err)
		return
	}

	fmt.Printf("resp: %#v. \n", resp)

	return
}
