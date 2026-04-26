package chat

// 消息类型
const (
	TypeText = 1 // 文本消息
)

// Message 客户端与服务端通信的消息结构体
type Message struct {
	Type     int    `json:"type"`      // 消息类型 1单聊 2群聊
	ToUserId int64  `json:"to_user_id"`// 接收者ID
	Content  string `json:"content"`   // 消息内容
	MsgType  int    `json:"msg_type"`  // 1文本
}