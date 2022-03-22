package message

// Text 文本消息
type Text struct {
	CommonToken
	Content CDATA `json:"Content"`
}

// NewText 初始化文本消息
func NewText(content string) *Text {
	text := new(Text)
	text.Content = CDATA(content)
	return text
}
