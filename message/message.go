package message

// MsgType 基本消息类型
type MsgType string

// EventType 事件类型
type EventType string

// InfoType 第三方平台授权事件类型
type InfoType string

const (
	// MsgTypeText 表示文本消息
	MsgTypeText MsgType = "text"
	// MsgTypeImage 表示图片消息
	MsgTypeImage = "image"
	// MsgTypeVoice 表示语音消息
	MsgTypeVoice = "voice"
	// MsgTypeVideo 表示视频消息
	MsgTypeVideo = "video"
	// MsgTypeMiniprogrampage 表示小程序卡片消息
	MsgTypeMiniprogrampage = "miniprogrampage"
	// MsgTypeShortVideo 表示短视频消息[限接收]
	MsgTypeShortVideo = "shortvideo"
	// MsgTypeLocation 表示坐标消息[限接收]
	MsgTypeLocation = "location"
	// MsgTypeLink 表示链接消息[限接收]
	MsgTypeLink = "link"
	// MsgTypeMusic 表示音乐消息[限回复]
	MsgTypeMusic = "music"
	// MsgTypeNews 表示图文消息[限回复]
	MsgTypeNews = "news"
	// MsgTypeTransfer 表示消息消息转发到客服
	MsgTypeTransfer = "transfer_customer_service"
	// MsgTypeEvent 表示事件推送消息
	MsgTypeEvent = "event"
)

const (
	// EventSubscribe 订阅
	EventSubscribe EventType = "subscribe"
	// EventUnsubscribe 取消订阅
	EventUnsubscribe = "unsubscribe"
	// EventScan 用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
	EventScan = "SCAN"
	// EventLocation 上报地理位置事件
	EventLocation = "LOCATION"
	// EventClick 点击菜单拉取消息时的事件推送
	EventClick = "CLICK"
	// EventView 点击菜单跳转链接时的事件推送
	EventView = "VIEW"
	// EventScancodePush 扫码推事件的事件推送
	EventScancodePush = "scancode_push"
	// EventScancodeWaitmsg 扫码推事件且弹出“消息接收中”提示框的事件推送
	EventScancodeWaitmsg = "scancode_waitmsg"
	// EventPicSysphoto 弹出系统拍照发图的事件推送
	EventPicSysphoto = "pic_sysphoto"
	// EventPicPhotoOrAlbum 弹出拍照或者相册发图的事件推送
	EventPicPhotoOrAlbum = "pic_photo_or_album"
	// EventPicWeixin 弹出微信相册发图器的事件推送
	EventPicWeixin = "pic_weixin"
	// EventLocationSelect 弹出地理位置选择器的事件推送
	EventLocationSelect = "location_select"
	// EventTemplateSendJobFinish 发送模板消息推送通知
	EventTemplateSendJobFinish = "TEMPLATESENDJOBFINISH"
	// EventMassSendJobFinish 群发消息推送通知
	EventMassSendJobFinish = "MASSSENDJOBFINISH"
	// EventWxaMediaCheck 异步校验图片/音频是否含有违法违规内容推送事件
	EventWxaMediaCheck = "wxa_media_check"
)

const (
	// 微信开放平台需要用到

	// InfoTypeVerifyTicket 返回ticket
	InfoTypeVerifyTicket InfoType = "component_verify_ticket"
	// InfoTypeAuthorized 授权
	InfoTypeAuthorized = "authorized"
	// InfoTypeUnauthorized 取消授权
	InfoTypeUnauthorized = "unauthorized"
	// InfoTypeUpdateAuthorized 更新授权
	InfoTypeUpdateAuthorized = "updateauthorized"
)

// MixMessage 存放所有微信发送过来的消息和事件
type MixMessage struct {
	CommonToken

	// 基本消息
	MsgID         int64   `json:"MsgId"` // 其他消息推送过来是MsgId
	TemplateMsgID int64   `json:"MsgID"` // 模板消息推送成功的消息是MsgID
	Content       string  `json:"Content"`
	Recognition   string  `json:"Recognition"`
	PicURL        string  `json:"PicUrl"`
	MediaID       string  `json:"MediaId"`
	Format        string  `json:"Format"`
	ThumbMediaID  string  `json:"ThumbMediaId"`
	LocationX     float64 `json:"Location_X"`
	LocationY     float64 `json:"Location_Y"`
	Scale         float64 `json:"Scale"`
	Label         string  `json:"Label"`
	Title         string  `json:"Title"`
	Description   string  `json:"Description"`
	URL           string  `json:"Url"`

	// 事件相关
	Event       EventType `json:"Event"`
	EventKey    string    `json:"EventKey"`
	Ticket      string    `json:"Ticket"`
	Latitude    string    `json:"Latitude"`
	Longitude   string    `json:"Longitude"`
	Precision   string    `json:"Precision"`
	MenuID      string    `json:"MenuId"`
	Status      string    `json:"Status"`
	SessionFrom string    `json:"SessionFrom"`
	TotalCount  int64     `json:"TotalCount"`
	FilterCount int64     `json:"FilterCount"`
	SentCount   int64     `json:"SentCount"`
	ErrorCount  int64     `json:"ErrorCount"`

	ScanCodeInfo struct {
		ScanType   string `json:"ScanType"`
		ScanResult string `json:"ScanResult"`
	} `json:"ScanCodeInfo"`

	SendPicsInfo struct {
		Count   int32      `json:"Count"`
		PicList []EventPic `json:"PicList>item"`
	} `json:"SendPicsInfo"`

	SendLocationInfo struct {
		LocationX float64 `json:"Location_X"`
		LocationY float64 `json:"Location_Y"`
		Scale     float64 `json:"Scale"`
		Label     string  `json:"Label"`
		Poiname   string  `json:"Poiname"`
	}

	// 第三方平台相关
	InfoType                     InfoType `json:"InfoType"`
	AppID                        string   `json:"AppId"`
	ComponentVerifyTicket        string   `json:"ComponentVerifyTicket"`
	AuthorizerAppid              string   `json:"AuthorizerAppid"`
	AuthorizationCode            string   `json:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int64    `json:"AuthorizationCodeExpiredTime"`
	PreAuthCode                  string   `json:"PreAuthCode"`

	// 卡券相关
	CardID              string `json:"CardId"`
	RefuseReason        string `json:"RefuseReason"`
	IsGiveByFriend      int32  `json:"IsGiveByFriend"`
	FriendUserName      string `json:"FriendUserName"`
	UserCardCode        string `json:"UserCardCode"`
	OldUserCardCode     string `json:"OldUserCardCode"`
	OuterStr            string `json:"OuterStr"`
	IsRestoreMemberCard int32  `json:"IsRestoreMemberCard"`
	UnionID             string `json:"UnionId"`

	// 内容审核相关
	IsRisky       bool   `json:"isrisky"`
	ExtraInfoJSON string `json:"extra_info_json"`
	TraceID       string `json:"trace_id"`
	StatusCode    int    `json:"status_code"`
}

// EventPic 发图事件推送
type EventPic struct {
	PicMd5Sum string `json:"PicMd5Sum"`
}

// ResponseEncryptedjsonMsg 需要返回的消息体
type ResponseEncryptedjsonMsg struct {
	jsonName     struct{} `json:"json" json:"-"`
	EncryptedMsg string   `json:"Encrypt"      json:"Encrypt"`
	MsgSignature string   `json:"MsgSignature" json:"MsgSignature"`
	Timestamp    int64    `json:"TimeStamp"    json:"TimeStamp"`
	Nonce        string   `json:"Nonce"        json:"Nonce"`
}

// CDATA  使用该类型,在序列化为 json 文本时文本会被解析器忽略
type CDATA string

// CommonToken 消息中通用的结构
type CommonToken struct {
	ToUserName   CDATA   `json:"ToUserName"`
	FromUserName CDATA   `json:"FromUserName"`
	CreateTime   int64   `json:"CreateTime"`
	MsgType      MsgType `json:"MsgType"`
}

// SetToUserName set ToUserName
func (msg *CommonToken) SetToUserName(toUserName CDATA) {
	msg.ToUserName = toUserName
}

// SetFromUserName set FromUserName
func (msg *CommonToken) SetFromUserName(fromUserName CDATA) {
	msg.FromUserName = fromUserName
}

// SetCreateTime set createTime
func (msg *CommonToken) SetCreateTime(createTime int64) {
	msg.CreateTime = createTime
}

// SetMsgType set MsgType
func (msg *CommonToken) SetMsgType(msgType MsgType) {
	msg.MsgType = msgType
}

// GetOpenID get the FromUserName value
func (msg *CommonToken) GetOpenID() string {
	return string(msg.FromUserName)
}
