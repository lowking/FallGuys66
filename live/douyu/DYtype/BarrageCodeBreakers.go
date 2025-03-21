package DYtype

import (
	"FallGuys66/live/douyu/lib"
	"FallGuys66/live/douyu/lib/logger"
	"bytes"
	"strings"
)

type InfoByteConfig struct {
	sendByte []byte
	endByte  []byte
}

type DyBarrageRawMsg struct {
	Self InfoByteConfig
	CodeBreakers
}

type CodeBreakers interface {
	Encode(msg string) []byte
	Decode(msg []byte) []string
	GetChatMessages(msgBytes []byte) []Response
	__parseMsg(rawMsg string) Response
}

/*
	type Response struct {
		Type  string `json:"type"`
		Rid   string `json:"rid"`
		Uid   string `json:"uid"`
		Nn    string `json:"nn"`
		Txt   string `json:"txt"`
		Cid   string `json:"cid"`
		Ic    string `json:"ic"`
		Level string `json:"level"`
		Sahf  string `json:"sahf"`
		Cst   string `json:"cst"`
		Bnn   string `json:"bnn"`
		Bl    string `json:"bl"`
		Brid  string `json:"brid"`
		Hc    string `json:"hc"`
		El    string `json:"el"`
		Lk    string `json:"lk"`
		Fl    string `json:"fl"`
		Urlev string `json:"urlev"`
		Dms   string `json:"dms"`
	}
*/
type Response map[string]string

type CodeBreakershandler struct {
	IsLive bool
}

// Encode:编码字符串
func (c CodeBreakershandler) Encode(msg string) []byte {
	dataLen := int32(len(msg) + 9)
	msgByte := []byte(msg)
	lenByte := lib.IntToBytes(dataLen, "little")
	sendByte := []byte{0xb1, 0x02, 0x00, 0x00}
	endByte := []byte{0x00}
	data := bytes.Join([][]byte{lenByte, lenByte, sendByte, msgByte, endByte}, []byte(""))
	return data
}

// Decode:解析返回字节码为字符串
func (c CodeBreakershandler) Decode(msgBytes []byte) []string {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("%v msgBytes：%x msgStr: %s", err, msgBytes, string(msgBytes))
		}
	}()
	pos := 0
	var msg []string
	for pos < len(msgBytes) {
		contentLength := lib.BytesToInt(msgBytes[pos:pos+4], "little")
		content := bytes.NewBuffer(msgBytes[pos+12 : pos+3+contentLength]).String()
		msg = append(msg, content)
		pos = 4 + contentLength + pos
	}
	return msg
}

// GetChatMessages:获取返回数据区分弹幕和其他信息
func (c CodeBreakershandler) GetChatMessages(msgByte []byte) []Response {
	decodeMsg := c.Decode(msgByte)
	var messages []Response
	for _, msg := range decodeMsg {
		res := c.__parseMsg(msg)
		messages = append(messages, res)
	}
	return messages
}

// __parseMsg:根据符号分割字符串组成消息
func (c CodeBreakershandler) __parseMsg(rawMsg string) Response {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("%v rawMsg：%s", err, rawMsg)
		}
	}()
	res := make(Response)
	attrs := strings.Split(rawMsg, "/")
	attrs = attrs[0:]
	for _, attr := range attrs {
		if attr != "" {
			attr := strings.Replace(attr, "@s", "/", 1)
			attr = strings.Replace(attr, "@A", "@", 1)
			couple := strings.Split(attr, "@=")

			if len(couple) < 2 {
				res[couple[0]] = ""
			} else {
				res[couple[0]] = couple[1]
			}
		}
	}
	return res
}
