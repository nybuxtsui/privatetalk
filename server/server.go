package server

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nybuxtsui/log"
	"github.com/urfave/cli/v2"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	writeCh = make(chan WriteMessage, 5000)
)

type WriteMessage struct {
	Message *TalkMessage
	Conn    *websocket.Conn
}

type TalkMessage struct {
	Id         string `json:"id,omitempty"`
	SrcType    string `json:"src,omitempty"`
	Type       string `json:"type,omitempty"`
	SenderId   string `json:"senderid,omitempty"`
	SenderName string `json:"sendername,omitempty"`
	Message    string `json:"msg,omitempty"`
}

func writeHandler() {
	for msg := range writeCh {
		b, err := json.Marshal(msg.Message)
		if err != nil {
			log.Error("marshal msg failed: %v", err.Error())
			continue
		}
		err = msg.Conn.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			log.Error("write msg failed: %v", err.Error())
			continue
		}
		log.Info("send: %s", b)
	}
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	userId := genId()
	userName := r.URL.Query().Get("username")
	roomId := r.URL.Query().Get("roomid")

	if userName == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "userName empty")
		return
	}
	if roomId == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "roomId empty")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade failed: %s", err.Error())
		return
	}

	c := make(chan LobbyResp, 1)

	defer func() {
		lobby.C <- LobbyReq{
			Type:       "disconnect",
			RoomId:     roomId,
			SenderId:   userId,
			SenderName: userName,
			C:          c,
		}
		<-c
		conn.Close()
	}()

	lobby.C <- LobbyReq{
		Type:       "join",
		RoomId:     roomId,
		SenderId:   userId,
		SenderName: userName,

		Conn: conn,
		C:    c,
	}
	<-c
	log.Info("connected")
	writeCh <- WriteMessage{
		Message: &TalkMessage{
			Type: "ok",
		},
		Conn: conn,
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Info("read failed: %s", err.Error())
			return
		}
		var msg TalkMessage
		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Info("unmarshal failed: %s", err.Error())
			return
		}
		switch msg.Type {
		case "msg":
			lobby.C <- LobbyReq{
				Type:       msg.Type,
				RoomId:     roomId,
				SenderId:   userId,
				SenderName: userName,
				Message:    msg.Message,
				C:          c,
			}
			<-c
			writeCh <- WriteMessage{
				Message: &TalkMessage{
					Id:      msg.Id,
					SrcType: msg.Type,
					Type:    "ok",
				},
				Conn: conn,
			}
		default:
			log.Error("unknown type: %s", msg.Type)
			return
		}
	}
}

func Run(c *cli.Context) error {
	log.Info("server start %s", c.String("addr"))

	go writeHandler()
	go lobby.worker()

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/chat", wsHandle)

	return http.ListenAndServe(c.String("addr"), nil)
}
