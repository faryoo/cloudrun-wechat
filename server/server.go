package server

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/faryoo/cloudrun-wechat/context"
	"github.com/faryoo/cloudrun-wechat/message"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime/debug"

	"github.com/faryoo/cloudrun-wechat/util"
)

// Server struct
type Server struct {
	*context.Context
	Writer  http.ResponseWriter
	Request *http.Request

	skipValidate bool

	openID string

	messageHandler func(*message.MixMessage) *message.Reply

	RequestRawXMLMsg  []byte
	RequestMsg        *message.MixMessage
	ResponseRawXMLMsg []byte
	ResponseMsg       interface{}

	isSafeMode bool
	random     []byte
	nonce      string
	timestamp  int64
}

// NewServer init
func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context

	return srv
}

// SkipValidate set skip validate
func (srv *Server) SkipValidate(skip bool) {
	srv.skipValidate = skip
}

// Serve 处理微信的请求消息
func (srv *Server) Serve() error {
	if !srv.Validate() {
		log.Error("Validate Signature Failed.")
		return fmt.Errorf("请求校验失败")
	}

	echostr, exists := srv.GetQuery("echostr")
	if exists {
		srv.String(echostr)

		return nil
	}

	response, err := srv.handleRequest()
	if err != nil {
		return err
	}

	// debug print request msg
	log.Debugf("request msg =%s", string(srv.RequestRawXMLMsg))

	return srv.buildResponse(response)
}

// Validate 校验请求是否合法
func (srv *Server) Validate() bool {

	return true
}

// HandleRequest 处理微信的请求
func (srv *Server) handleRequest() (reply *message.Reply, err error) {
	// set isSafeMode
	srv.isSafeMode = false
	encryptType := srv.Query("encrypt_type")

	if encryptType == "aes" {
		srv.isSafeMode = true
	}

	// set openID
	srv.openID = srv.Query("openid")

	var msg interface{}
	msg, err = srv.getMessage()

	if err != nil {
		return
	}

	mixMessage, success := msg.(*message.MixMessage)

	if !success {
		err = errors.New("消息类型转换失败")
	}

	srv.RequestMsg = mixMessage
	reply = srv.messageHandler(mixMessage)

	return
}

// GetOpenID return openID
func (srv *Server) GetOpenID() string {
	return srv.openID
}

// getMessage 解析微信返回的消息
func (srv *Server) getMessage() (interface{}, error) {
	var rawXMLMsgBytes []byte

	var err error

	rawXMLMsgBytes, err = ioutil.ReadAll(srv.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("从body中解析xml失败, err=%v", err)
	}

	srv.RequestRawXMLMsg = rawXMLMsgBytes

	return srv.parseRequestMessage(rawXMLMsgBytes)
}

func (srv *Server) parseRequestMessage(rawXMLMsgBytes []byte) (msg *message.MixMessage, err error) {
	msg = &message.MixMessage{}
	err = xml.Unmarshal(rawXMLMsgBytes, msg)

	return
}

// SetMessageHandler 设置用户自定义的回调方法
func (srv *Server) SetMessageHandler(handler func(*message.MixMessage) *message.Reply) {
	srv.messageHandler = handler
}

func (srv *Server) buildResponse(reply *message.Reply) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic error: %v\n%s", e, debug.Stack())
		}
	}()
	if reply == nil {
		// do nothing
		return nil
	}
	msgType := reply.MsgType
	switch msgType {
	case message.MsgTypeText:
	case message.MsgTypeImage:
	case message.MsgTypeVoice:
	case message.MsgTypeVideo:
	case message.MsgTypeMusic:
	case message.MsgTypeNews:
	case message.MsgTypeTransfer:
	default:
		err = message.ErrUnsupportReply
		return
	}

	msgData := reply.MsgData
	value := reflect.ValueOf(msgData)
	// msgData must be a ptr
	kind := value.Kind().String()
	if kind != "ptr" {
		return message.ErrUnsupportReply
	}

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(srv.RequestMsg.FromUserName)
	value.MethodByName("SetToUserName").Call(params)

	params[0] = reflect.ValueOf(srv.RequestMsg.ToUserName)
	value.MethodByName("SetFromUserName").Call(params)

	params[0] = reflect.ValueOf(msgType)
	value.MethodByName("SetMsgType").Call(params)

	params[0] = reflect.ValueOf(util.GetCurrTS())
	value.MethodByName("SetCreateTime").Call(params)

	srv.ResponseMsg = msgData
	srv.ResponseRawXMLMsg, err = xml.Marshal(msgData)
	return
}

// Send 将自定义的消息发送
func (srv *Server) Send() (err error) {
	replyMsg := srv.ResponseMsg
	log.Debugf("response msg =%+v", replyMsg)

	if replyMsg != nil {
		srv.XML(replyMsg)
	}

	return
}
