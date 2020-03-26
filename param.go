package wechatpay

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
)

type Params map[string]string

func (p Params) Get(key string) string {
	v, _ := p[key]

	return v
}

func (p Params) Set(key string, value interface{}) {
	if value == nil {
		return
	}

	switch value.(type) {
	case int8, int16, int32, int64:
		v, _ := value.(int64)
		s := strconv.FormatInt(v, 10)
		p[key] = s
		break

	case uint8, uint16, uint32, uint64:
		v, _ := value.(uint64)
		s := strconv.FormatUint(v, 10)
		p[key] = s
		break

	case float32, float64:
		v, _ := value.(float64)
		s := strconv.FormatFloat(v, 'f', 0, 64)
		p[key] = s
		break

	case bool:
		v, _ := value.(bool)
		s := strconv.FormatBool(v)
		p[key] = s
		break

	case string:
		v, _ := value.(string)
		p[key] = v
		break
	}

	return
}

func (p Params) SetInt(key string, value int64) {
	s := strconv.FormatInt(value, 10)
	p[key] = s

	return
}

func (p Params) SetUint(key string, value uint64) {
	s := strconv.FormatUint(value, 10)
	p[key] = s

	return
}

func (p Params) SetFloat(key string, value float64) {
	s := strconv.FormatFloat(value, 'f', 0, 64)
	p[key] = s

	return
}

func (p Params) SetBool(key string, value bool) {
	s := strconv.FormatBool(value)
	p[key] = s

	return
}

func (p Params) SetString(key string, value string) {
	p[key] = value

	return
}

// 格式化 map[string]string 为 xml 格式, 根节点名字为 xml.
func (p Params) FormatParams2XML(xmlWriter io.Writer) (err error) {
	if xmlWriter == nil {
		return errors.New("nil xmlWriter")
	}
	if _, err = io.WriteString(xmlWriter, "<xml>"); err != nil {
		return
	}

	for k, v := range p {
		if _, err = io.WriteString(xmlWriter, "<"+k+">"); err != nil {
			return
		}
		if err = xml.EscapeText(xmlWriter, []byte(v)); err != nil {
			return
		}
		if _, err = io.WriteString(xmlWriter, "</"+k+">"); err != nil {
			return
		}
	}

	if _, err = io.WriteString(xmlWriter, "</xml>"); err != nil {
		return
	}
	return
}

// 解析xml, 返回第一级子节点的键值集合, 如果第一级子节点包含有子节点, 则跳过.
func ParseXML2Params(xmlReader io.Reader) (p Params, err error) {
	if xmlReader == nil {
		err = errors.New("nil xmlReader")
		return
	}
	d := xml.NewDecoder(xmlReader)
	var key string          // 当前"第一级"子节点的 key
	var buffer bytes.Buffer // 当前"第一级"子节点的 value
	var depth int           // 当前节点的深度

	p = make(Params)
	for {
		var tk xml.Token
		tk, err = d.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
				return
			}
			return
		}

		switch v := tk.(type) {
		case xml.StartElement:
			depth++
			switch depth {
			case 1: // do nothing
			case 2:
				key = v.Name.Local
				buffer.Reset()
			case 3:
				if err = d.Skip(); err != nil {
					return
				}
				depth--
				key = "" // key == "" 暗示了当前第一级子节点包含子节点
			default:
				panic("incorrect algorithm")
			}
		case xml.CharData:
			if depth == 2 && key != "" {
				buffer.Write(v)
			}
		case xml.EndElement:
			if depth == 2 && key != "" {
				p[key] = buffer.String()
			}
			depth--
		}
	}

	return
}
