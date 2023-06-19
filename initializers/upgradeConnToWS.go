package initializers

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:	  func(r *http.Request) bool {
			return true
		},
	}
	ChannelMap = make(map[string]*websocket.Conn)
	ChannelMutex sync.Mutex
)