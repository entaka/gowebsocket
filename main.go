package main

import (
	"flag"
	"fmt"
	"gowebsocket/modules/multi"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// コマンドで入力された引数を処理する
// .String() 文字列
// .Int() 数値
// .Bool() 成否
// ("取得する引数","デフォルト値","表示するメッセージの設定")
var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html") // 静的ファイルの指定
}
func serveMulti(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "multi.html")
}

func main() {
	flag.Parse()    // 実行引数を設定しておいた変数に挿入する
	hub := newHub() // importしていないが、同時にrunさせることで別ファイル関数を呼び出せる 同時実行:hub.go
	go hub.run()    // hubを使いclientと同期をとる

	// websocket処理を分ける場合はhubを増やす
	// 同時に処理を行うものでは同じhubを使用することで共通化できる
	hub2 := newHub() // importしていないが、同時にrunさせることで別ファイル関数を呼び出せる 同時実行:hub.go
	go hub2.run()    // hubを使いclientと同期をとる

	http.HandleFunc("/", serveHome) // 表示ページ用処理
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r) // 同時実行:client.go
	})
	http.HandleFunc("/ws2", func(w http.ResponseWriter, r *http.Request) {
		serveWs2(hub2, w, r) // 同時実行:client.go
	})
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(r.RequestURI)
		query := u.Query()
		// query -> map[a:[AAA] b:[BBB] c:[CCC] d:[DDD]]
		i, _ := strconv.Atoi(query["id"][0])
		id := auth2(i)
		fmt.Fprint(w, id)
	})

	// アクセス毎にwebsocketを作る
	// 後々にwebsocketの数をコントロールしいくつかのアクセスごとにwebsocketを作る
	http.HandleFunc("/multiws", func(w http.ResponseWriter, r *http.Request) {
		multi.ServeWsMulti(w, r)
	})
	http.HandleFunc("/multi", serveMulti) // 表示ページ用処理

	http.Handle("/js/", http.FileServer(http.Dir("assets/")))
	err := http.ListenAndServe(*addr, nil) // 指定ポートでサーバーを立てる
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
