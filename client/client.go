package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	"github.com/oyvindsk/go-protobuf-ws-chat/lib/message"
)

func main() {

	msgtype := websocket.BinaryMessage

	var ws *websocket.Conn

	curState := stateDisconnected
	curEvent := eventConnect
	var exitAction action
	for {
		var err error

		curState, exitAction, err = curState.Next(curEvent)
		if err != nil {
			log.Fatal(err)
		}
		curEvent++

		log.Println(exitAction)

		switch event(exitAction) {

		case eventConnect:
			log.Println("Action connect")
			dialer := websocket.Dialer{}
			ws, _, err = dialer.Dial("ws://localhost:8080/ws", http.Header{})
			if err != nil {
				log.Fatal("Error connecting:", err)
			}

			log.Println("Connected!")

		case eventHandshake:
			// Expect this to go like this:
			// >> YO! 1
			// << OY! 1
			// >> OK
			log.Println("Action handshake")
			wsMustWriteStr(ws, msgtype, "YO! 1")
			wsMustReadStr(ws, msgtype, "OY! 1")
			wsMustWriteStr(ws, msgtype, "OK")

		case eventRegister:
			log.Println("Action Register")

			reg := message.RegisterNick{}
			reg.Nick = "OleP"
			regBytes, err := proto.Marshal(&reg)
			if err != nil {
				log.Fatal(err)
			}

			msg := message.Message{}
			msg.Type = message.MessageType_REGISTERNICK
			msg.Content = regBytes
			msgBytes, err := proto.Marshal(&msg)
			if err != nil {
				log.Fatal(err)
			}

			wsMustWrite(ws, msgtype, msgBytes)

		case eventDisconnect:
			log.Println("Action Disconnect")

		case eventJoinRoom:
			log.Println("Action Join Room")

		case eventLeaveRoom:
			log.Println("Action Leave Room")

		case eventChatMsg:
			log.Println("Action Chat Message")

		}

	}

	log.Fatal("done")

	// Nice! Switch to proto3
	msg := message.Message{}
	msg.Type = message.MessageType_CHATMESSAGE

	chat := message.ChatMessage{}
	chat.From = "Ole"
	chat.To = "Petter"
	chat.Data = "hallo hallo hallo =)"

	chatBytes, err := proto.Marshal(&chat)
	if err != nil {
		log.Fatal(err)
	}

	msg.Content = chatBytes
	data, err := proto.Marshal(&msg)

	wsMustWrite(ws, msgtype, data)

}

// Disconnected    Connected       Handshaked      Registered
const (
	stateDisconnected = state(iota)
	stateConnected
	stateHandshaked
	stateRegistered
)

const (
	eventConnect = event(iota)
	eventHandshake
	eventRegister
	eventDisconnect
	eventJoinRoom
	eventLeaveRoom
	eventChatMsg
)

// FIXME
type state int
type event int
type action int

// returns the next state and the "action" that has to be done to get there
func (s state) Next(e event) (state, action, error) {
	switch s {

	case stateDisconnected:
		switch e {

		case eventConnect:
			// do something usefull
			return stateConnected, action(e), nil
		}

	case stateConnected:
		switch e {

		case eventHandshake:
			return stateHandshaked, action(e), nil

		case eventDisconnect:
			return stateDisconnected, action(e), nil
		}

	case stateHandshaked:
		switch e {

		case eventRegister:
			return stateRegistered, action(e), nil

		case eventDisconnect:
			return stateDisconnected, action(e), nil
		}

	case stateRegistered:
		switch e {

		case eventRegister:
			return stateRegistered, action(e), nil

		case eventJoinRoom:
			return stateRegistered, action(e), nil

		case eventLeaveRoom:
			return stateRegistered, action(e), nil

		case eventChatMsg:
			return stateRegistered, action(e), nil

		case eventDisconnect:
			return stateDisconnected, action(e), nil

		}
	}
	return -1, -1, fmt.Errorf("%s", "Unsupported state or event")
}

//type stateFunc func() stateFunc

// do whatever is necessary to change to the next state, and return that next state

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
