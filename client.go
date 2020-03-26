package wechatpay

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type Client struct {
	stdClient *http.Client

	appId  string
	apiKey string
	mchId  string
}

func NewClient(appId, apiKey, mchId string) *Client {
	return &Client{
		stdClient: &http.Client{},
		appId:     appId,
		apiKey:    apiKey,
		mchId:     mchId,
	}
}

func (c *Client) withCert(certFile, keyFile, rootFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(rootFile)
	if err != nil {
		return err
	}

	certPoolFunc := x509.NewCertPool()

	ok := certPoolFunc.AppendCertsFromPEM(data)
	if !ok {
		return errors.New("parse PEM encoded certificates failed")
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPoolFunc,
	}

	transport := &http.Transport{
		TLSClientConfig: conf,
	}

	if c.stdClient != nil {
		c.stdClient.Transport = transport
	} else {
		c.stdClient = &http.Client{
			Transport: transport,
		}
	}

	return nil
}

func (c *Client) SetAppId(appId string) error {
	if c == nil {
		return errors.New("client can't be nil")
	}

	appId = strings.TrimSpace(appId)
	if appId == "" {
		return errors.New("appId can't be nil")
	}

	c.appId = appId

	return nil
}

func (c *Client) SetApiKey(apiKey string) error {
	if c == nil {
		return errors.New("client can't be nil")
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return errors.New("apiKey can't be nil")
	}

	c.apiKey = apiKey

	return nil
}

func (c *Client) SetMchId(mchId string) error {
	if c == nil {
		return errors.New("client can't be nil")
	}

	mchId = strings.TrimSpace(mchId)
	if mchId == "" {
		return errors.New("mchId can't be nil")
	}

	c.mchId = mchId

	return nil
}

func (c *Client) Post(url string, contentType string, body io.Reader) (resp Params, err error) {
	if body == nil {
		return nil, errors.New("body can't be nil")
	}

	httpResp, err := c.stdClient.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("StatusCode %v", httpResp.StatusCode)
		return nil, err
	}

	return ParseXML2Params(httpResp.Body)
}

//校验签名
func (c *Client) CheckSign(ps Params) bool {
	return ps.Get("sign") == c.Sign(ps, nil)
}

//生成签名，默认使用MD5
func (c *Client) Sign(p Params, fn func() hash.Hash) (sign string) {
	ks := make([]string, 0, len(p))
	for k := range p {
		if k == "sign" {
			continue
		}
		ks = append(ks, k)
	}
	sort.Strings(ks)

	if fn == nil {
		fn = md5.New
	}
	h := fn()
	signature := make([]byte, h.Size()*2)

	for _, k := range ks {
		v := p[k]
		if v == "" {
			continue
		}
		h.Write([]byte(k))
		h.Write([]byte{'='})
		h.Write([]byte(v))
		h.Write([]byte{'&'})
	}
	h.Write([]byte("key="))
	h.Write([]byte(c.apiKey))

	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}
