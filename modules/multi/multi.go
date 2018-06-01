package multi

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second // 通信タイムアウト
	pongWait       = 60 * time.Second // 接続開始待機時間
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 // 最大受け取り文字数
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ws client
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func auth() int {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(1000)
	return id
}

// 以下複数websocketテスト
func ServeWsMulti(w http.ResponseWriter, r *http.Request) {
	log.Printf("ServeWsMulti: %s", "new Hub test")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade: %s", "upgrade err")
		log.Println(err)
		return
	}
	log.Printf("NewHub: %s", "NewHub")
	hub := NewHub()
	go hub.Run()
	log.Printf("NewHub: %s", "New client")
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	log.Printf("NewHub: %s", "New hub register")
	client.hub.register <- client
	log.Printf("register: %s", "client")
	go client.writePumpMulti()
	go client.readPumpMulti()
}

func newHubMulti() {

}

// 受け取り処理
// 永久ループでメッセージを待ち続ける
func (c *Client) readPumpMulti() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("message: %s", message)
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// 送信処理
func (c *Client) writePumpMulti() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			log.Printf("send_len: %d", n)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			w.Write([]byte("------------------------------------------"))
			log.Printf("message: %s", message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
