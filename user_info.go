package wechatpay

import (
	"io/ioutil"
	"net/http"
)

/*从微信服务器拉取玩家的微信信息*/
func GetUserInfoByAccessTokenAndOpenId(accessToken string, openId string) ([]byte, error) {
	//userInfoURL格式： https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID

	userInfoURL := "https://api.weixin.qq.com/sns/userinfo?" + "access_token=" + accessToken + "&openid=" + openId

	var err error

	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
