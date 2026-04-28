package chat

// 消息类型
const (
	TypeText int8 = 1 // 文本消息

	TypePrivate int8 = 1 // 单聊
	TypeGroup   int8 = 2 // 群聊
)

// Message 客户端与服务端通信的消息结构体
type Message struct {
	Type     int8    `json:"type"`      // 消息类型 1单聊 2群聊
	ToUserId int64  `json:"to_user_id"`// 接收者ID
	Content  string `json:"content"`   // 消息内容
	MsgType  int8    `json:"msg_type"`  // 1文本
}