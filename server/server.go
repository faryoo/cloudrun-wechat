package server

import (
	"errors"
	"fmt"
	"github.com/faryoo/cloudrun-wechat/context"
	"github.com/faryoo/cloudrun-wechat/message"
	"github.com/faryoo/cloudrun-wechat/util"
	json "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime/debug"
)

type Server struct {
	*context.Context
	Writer       http.ResponseWriter
	Request      *http.Request
	skipValidate bool
	openID       string

	messageHandler func(*message.MixMessage) *message.Reply

	RequestRawJSONMsg  []byte
	RequestMsg         *message.MixMessage
	ResponseRawJSONMsg []byte
	ResponseMsg        interface{}

	random    []byte
	timestamp int64
}

func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context

	return srv
}
func (srv *Server) SkipValidate(skip bool) {
	srv.skipValidate = skip
}
func (srv *Server) Serve() error {
	if !srv.Validate() {
		log.Error("Validate Signature Failed.")
		return fmt.Errorf("请求校验失败")
	}

	response, err := srv.handleRequest()
	if err != nil {
		return err
	}

	// debug print request msg
	log.Debugf("request msg =%s", string(srv.RequestRawJSONMsg))

	return srv.buildResponse(response)
}
func (srv *Server) Validate() bool {
	if srv.skipValidate {
		return true
	}

	return true
}
func (srv *Server) handleRequest() (reply *message.Reply, err error) {
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
func (srv *Server) GetOpenID() string {
	return srv.openID
}
func (srv *Server) getMessage() (interface{}, error) {
	var rawJSONMsgBytes []byte

	var err error

	rawJSONMsgBytes, err = ioutil.ReadAll(srv.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("从body中解析json失败, err=%v", err)
	}

	srv.RequestRawJSONMsg = rawJSONMsgBytes

	return srv.parseRequestMessage(rawJSONMsgBytes)
}

func (srv *Server) parseRequestMessage(rawJSONMsgBytes []byte) (msg *message.MixMessage, err error) {
	msg = &message.MixMessage{}
	err = json.Unmarshal(rawJSONMsgBytes, msg)

	return
}
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
	srv.ResponseRawJSONMsg, err = json.Marshal(msgData)
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
