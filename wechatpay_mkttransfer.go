package wechatpay

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	URL_MKTTRANSFER = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	CONTENT_TYPE    = "application/xml; charset=utf-8"
)

func (c *Client) MktTransfer(ps Params, certFile, keyFile, rootFile string) (p Params, err error) {
	if ps == nil {
		return nil, errors.New("input Params can't be nil")
	}

	certFile = strings.TrimSpace(certFile)
	keyFile = strings.TrimSpace(keyFile)
	rootFile = strings.TrimSpace(rootFile)

	if certFile == "" || keyFile == "" || rootFile == "" {
		return nil, errors.New("certFile, keyFile or rootFile can't be nil")
	}

	//检查必要参数
	//必要参数参考https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_2
	requisiteParams := []string{"mch_appid", "mchid", "nonce_str", "sign", "partner_trade_no",
		"openid", "check_name", "amount", "desc", "spbill_create_ip"}

	if _, ok := ps["mch_appid"]; !ok {
		ps["mch_appid"] = c.appId
	}

	if _, ok := ps["mchid"]; !ok {
		ps["mchid"] = c.mchId
	}

	if _, ok := ps["nonce_str"]; !ok {
		ps["nonce_str"] = UUID(16)
	}

	switch ps["check_name"] {
	case "NO_CHECK", "FORCE_CHECK":
		break

	default:
		ps["check_name"] = "NO_CHECK"
		break
	}

	ps["sign"] = ""

	for _, v := range requisiteParams {
		if _, ok := ps[v]; !ok {
			return nil, fmt.Errorf("lost requisite params %v", v)
		}
	}

	err = c.withCert(certFile, keyFile, rootFile)
	if err != nil {
		return nil, err
	}

	ps["sign"] = c.Sign(ps, nil)

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	if err := ps.FormatParams2XML(buf); err != nil {
		return nil, err
	}

	p, err = c.Post(URL_MKTTRANSFER, CONTENT_TYPE, buf)
	if err != nil {
		return nil, err
	}

	return handleError(p)
}

func handleError(ps Params) (p Params, err error) {
	if ps == nil {
		return nil, errors.New("input Params can't be nil")
	}

	switch ps["return_code"] {
	case "FAIL":
		p = nil
		err = fmt.Errorf("%v", ps["return_msg"])
		break

	case "SUCCESS":
		if ps["result_code"] == "SUCCESS" {
			p = ps
			err = nil

			break
		}

		p = nil
		err = fmt.Errorf("%v", ps["err_code_des"])
		break

	default:
		p = nil
		err = fmt.Errorf("invaild return_code %v. ", ps["return_code"])
	}

	return
}
