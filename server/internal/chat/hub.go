package chat 


import (
	"go.uber.org/zap"
	"sync"
)

type Hub struct {
	// 在线用户列表：key=userID, value=*Client
	OnlineUsers map[int64]*Client 
	mutex sync.RWMutex 
}

// 全局单例
var H = &Hub {
	OnlineUsers: make(map[int64]*Client),
}

// 注册连接（用户上线）
func (h *Hub) Register(userID int64, c *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.OnlineUsers[userID] = c
	zap.S().Infof("用户[%d]上线，当前在线：%d", userID, len(h.OnlineUsers))
}

// 注销连接（用户下线）
func (h *Hub) UnRegister(userID int64) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.OnlineUsers, userID)
	zap.S().Infof("用户[%d]下线，当前在线：%d", userID, len(h.OnlineUsers))
}

// 获取用户在线状态
func (h *Hub) GetClient(userID int64) (*Client, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	client, ok := h.OnlineUsers[userID]
	return client, ok 
}