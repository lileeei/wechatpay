package wechatpay

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	URL_UNIFIEDORDER = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

func (c *Client) UnifiedOrderJSAPI(ps Params) (p Params, err error) {
	if ps == nil {
		return nil, errors.New("input Params can't be nil")
	}

	//检查参数
	//必要参数参考https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
	requisiteParams := []string{"appid", "mch_id", "nonce_str", "sign", "body", "out_trade_no", "total_fee", "spbill_create_ip", "notify_url", "trade_type", "openid"}

	ps["trade_type"] = "JSAPI"

	if _, ok := ps["appid"]; !ok {
		ps["appid"] = c.appId
	}

	if _, ok := ps["mch_id"]; !ok {
		ps["mch_id"] = c.mchId
	}

	if _, ok := ps["nonce_str"]; !ok {
		ps["nonce_str"] = UUID(16)
	}

	ps["sign"] = ""

	for _, v := range requisiteParams {
		if _, ok := ps[v]; !ok {
			return nil, fmt.Errorf("lost requisite params %v", v)
		}
	}

	ps["sign"] = c.Sign(ps, nil)

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	if err := ps.FormatParams2XML(buf); err != nil {
		return nil, err
	}

	httpResp, err := c.Post(URL_UNIFIEDORDER, CONTENT_TYPE, buf)
	if err != nil {
		return nil, err
	}

	respParams, err := handleError(httpResp)
	if err != nil {
		return nil, err
	}

	if !c.CheckSign(respParams) {
		return nil, errors.New("check signature failed")
	}

	p = make(Params)

	p["appId"] = c.appId
	p["partnerId"] = c.mchId
	p["package"] = fmt.Sprintf("prepay_id=%v", respParams["prepay_id"])
	p["nonceStr"] = respParams["nonce_str"]
	p["timeStamp"] = strconv.Itoa(int(time.Now().Unix()))
	p["signType"] = "MD5"

	p["paySign"] = c.Sign(p, nil)

	return p, nil
}

func (c *Client) UnifiedOrderAPP(ps Params) (p Params, err error) {
	if ps == nil {
		return nil, errors.New("input Params can't be nil")
	}

	//检查参数
	//必要参数参考https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
	requisiteParams := []string{"appid", "mch_id", "nonce_str", "sign", "body", "out_trade_no", "total_fee", "notify_url", "trade_type"}

	ps["trade_type"] = "APP"

	if _, ok := ps["appid"]; !ok {
		ps["appid"] = c.appId
	}

	if _, ok := ps["mchid"]; !ok {
		ps["mchid"] = c.mchId
	}

	if _, ok := ps["nonce_str"]; !ok {
		ps["nonce_str"] = UUID(16)
	}

	for _, v := range requisiteParams {
		if _, ok := ps[v]; !ok {
			return nil, fmt.Errorf("lost requisite params %v", v)
		}
	}

	ps["sign"] = c.Sign(ps, nil)

	buf := bytes.NewBuffer(make([]byte, 0, 16<<10))
	if err := ps.FormatParams2XML(buf); err != nil {
		return nil, err
	}

	httpResp, err := c.Post(URL_UNIFIEDORDER, CONTENT_TYPE, buf)
	if err != nil {
		return nil, err
	}

	respParams, err := handleError(httpResp)
	if err != nil {
		return nil, err
	}

	if !c.CheckSign(respParams) {
		return nil, errors.New("check signature failed")
	}

	p = make(Params)

	p["appId"] = c.appId
	p["partnerid"] = c.mchId
	p["package"] = "Sign=WXPay"
	p["nonceStr"] = respParams["nonce_str"]
	p["timeStamp"] = strconv.Itoa(int(time.Now().Unix()))
	p["extData"] = "app data"
	p["sign"] = c.Sign(p, nil)

	return p, nil
}
