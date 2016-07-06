package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	"github.com/oyvindsk/go-protobuf-ws-chat/lib/message"
)

func main() {

	msgtype := websocket.BinaryMessage

	dialer := websocket.Dialer{}
	ws, _, err := dialer.Dial("ws://localhost:8080/ws", http.Header{})
	if err != nil {
		log.Println("Error connecting:", err)
		return
	}

	log.Println("Connected!")

	// Expect this to go like this:
	// >> YO! 1
	// << OY! 1
	// >> OK

	wsMustWriteStr(ws, msgtype, "YO! 1")
	wsMustReadStr(ws, msgtype, "OY! 1")
	wsMustWriteStr(ws, msgtype, "OK")

	// Nice! Switch to proto3
	msg := message.Message{}
	msg.From = "Ole"
	msg.Foo = 44
	data, err := proto.Marshal(&msg)
	wsMustWrite(ws, msgtype, data)

}

func wsMustWrite(ws *websocket.Conn, msgtype int, data []byte) {
	if err := ws.WriteMessage(msgtype, data); err != nil {
		log.Fatal("Ws Writing failed:", err)
	}
}

func wsMustWriteStr(ws *websocket.Conn, msgtype int, data string) {
	wsMustWrite(ws, msgtype, []byte(data))
}

func wsMustRead(ws *websocket.Conn, expectedType int, expectedData []byte) {
	msgtype, msg, err := ws.ReadMessage()
	if err != nil {
		log.Fatal("Ws Reading failed:", err)
	}
	if msgtype != expectedType {
		log.Fatalf("Ws Reading failed. Got type: %d, expected: %d", msgtype, expectedType)
	}
	if !bytes.Equal(msg, expectedData) {
		log.Fatalf("Ws Reading failed. Got data: %s, expected: %s", msg, expectedData)
	}

}

func wsMustReadStr(ws *websocket.Conn, expectedType int, expectedData string) {
	wsMustRead(ws, expectedType, []byte(expectedData))
}
