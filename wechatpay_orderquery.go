package wechatpay

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	URL_ORDERQUERY = "https://api.mch.weixin.qq.com/pay/orderquery"
)

func (c *Client) OrderQuery(ps Params) (p Params, err error) {
	if ps == nil {
		return nil, errors.New("input Params can't be nil")
	}

	//检查必要参数
	//必要参数参考https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
	requisiteParams := []string{"appid", "mch_id", "out_trade_no", "nonce_str", "sign", "sign_type"}

	ps["appid"] = c.appId
	ps["mch_id"] = c.mchId

	ps["sign"] = ""

	if _, ok := ps["nonce_str"]; !ok {
		ps["nonce_str"] = UUID(16)
	}

	if _, ok := ps["sign_type"]; !ok {
		ps["sign_type"] = "MD5"
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

	httpResp, err := c.Post(URL_ORDERQUERY, CONTENT_TYPE, buf)
	if err != nil {
		return nil, err
	}

	respParams, err := handleError(httpResp)
	if err != nil {
		return nil, err
	}

	return respParams, nil
}
