package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// box point
type Box struct {
	PosX   int8 `json:"x"`
	PosY   int8 `json:"y"`
	Width  int8 `json:"w"`
	Height int8 `json:"h"`
}

var box Box

type User struct {
	UserID  int `json:"user_id"`
	UserBox Box `json:"user_box"`
}

var UserMaps map[int]User = make(map[int]User)

func auth2(id int) int {
	b := Box{}
	b.PosX = 100
	b.PosY = 100
	b.Width = 100
	b.Height = 100
	u := User{}
	u.UserID = id
	u.UserBox = b
	log.Printf("box: %+v", b)
	log.Printf("u: %+v", u)
	UserMaps[id] = u
	log.Printf("UserMaps : %+v", UserMaps)
	return id
}

// 以下複数websocketテスト
func serveWs2(hub *Hub, w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.RequestURI)
	query := u.Query()
	i, _ := strconv.Atoi(query["id"][0])

	log.Printf("Request hub: %+v", hub)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), id: i}
	client.hub.register <- client

	go client.writePump2()
	go client.readPump2()
}

// 受け取り処理
// 永久ループでメッセージを待ち続ける
func (c *Client) readPump2() {
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
func (c *Client) writePump2() {
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

			if val, ok := UserMaps[c.id]; ok {
				log.Printf("login: %+v", val)
				switch fmt.Sprintf("%s", message) {
				case "ArrowLeft":
					val.UserBox.PosX -= 5
				case "ArrowRight":
					val.UserBox.PosX += 5
				case "ArrowUp":
					val.UserBox.PosY -= 5
				case "ArrowDown":
					val.UserBox.PosY += 5
				default:
					log.Printf("no UserMapMessage Key")
				}
			} else {
				log.Printf("no UserMaps: ")
			}

			log.Printf("send: %s", message)
			log.Printf("box send: %+v", box)
			log.Printf("UserMaps: %+v", UserMaps)
			// j, _ := json.Marshal(box)
			j, _ := json.Marshal(UserMaps)
			log.Printf("send: %s", j)
			//str := fmt.Sprintf("%+v", box)
			w.Write([]byte(j))

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
