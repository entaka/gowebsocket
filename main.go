package main

import (
	"flag"
	"log"
	"net/http"
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

func main() {
	flag.Parse()
	hub := newHub() // importしていないが、同時にrunさせることで別ファイル関数を呼び出せる 同時実行:hub.go
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r) // 同時実行:client.go
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
