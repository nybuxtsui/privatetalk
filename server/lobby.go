package server

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/nybuxtsui/log"
)

type LobbyReq struct {
	Type string

	RoomId   string
	RoomName string

	SenderId   string
	SenderName string

	Message string

	Conn *websocket.Conn
	C    chan LobbyResp
}

type LobbyResp struct {
	Status string
	RoomId string
}

type Lobby struct {
	Rooms map[string]*Room

	C chan LobbyReq
}

func (lobby *Lobby) worker() {
	for msg := range lobby.C {
		switch msg.Type {
		case "msg":
			room, ok := lobby.Rooms[msg.RoomId]
			if !ok {
				log.Error("room not found: %v", msg.RoomId)
				if msg.C != nil {
					msg.C <- LobbyResp{Status: "room_not_found"}
				}
				continue
			}
			room.sendMsg(msg.Type, msg.SenderId, msg.SenderName, msg.Message)
			if msg.C != nil {
				msg.C <- LobbyResp{
					Status: "ok",
				}
			}
		case "create-room":
			room := &Room{
				Id:    genId(),
				Name:  msg.RoomName,
				Users: make(map[string]*User),
			}
			lobby.Rooms[room.Id] = room
			if msg.C != nil {
				msg.C <- LobbyResp{
					Status: "ok",
					RoomId: room.Id,
				}
			}
		case "join":
			room, ok := lobby.Rooms[msg.RoomId]
			if !ok {
				room = &Room{
					Id:    msg.RoomId,
					Name:  msg.RoomName,
					Users: make(map[string]*User),
				}
				log.Info("create room: %v", msg.RoomId)
				lobby.Rooms[room.Id] = room
			}
			room.Users[msg.SenderId] = &User{
				Id:   msg.SenderId,
				Name: msg.SenderName,
				Conn: msg.Conn,
			}
			room.sendMsg(
				"join",
				msg.SenderId,
				msg.SenderName,
				"",
			)
			if msg.C != nil {
				msg.C <- LobbyResp{
					Status: "ok",
					RoomId: room.Id,
				}
			}
		case "disconnect":
			room, ok := lobby.Rooms[msg.RoomId]
			if ok {
				delete(room.Users, msg.SenderId)
				if len(room.Users) == 0 {
					delete(lobby.Rooms, room.Id)
				} else {
					room.sendMsg(
						"leave",
						"",
						"system",
						fmt.Sprintf("%v(%v)离开聊天室", msg.SenderName, msg.SenderId),
					)
				}
			}
		}
	}
}

type Room struct {
	Id    string
	Name  string
	Users map[string]*User
}

func (room *Room) sendMsg(msgType, id, name, msg string) {
	for keyid, user := range room.Users {
		if keyid != id {
			user.sendMsg(msgType, id, name, msg)
		}
	}
}

type User struct {
	Id   string
	Name string

	Conn *websocket.Conn
}

func (user *User) sendMsg(msgType, id, name, msg string) {
	writeCh <- WriteMessage{
		Message: &TalkMessage{
			Type:       msgType,
			SenderId:   id,
			SenderName: name,
			Message:    msg,
		},
		Conn: user.Conn,
	}
}

var lobby = Lobby{
	make(map[string]*Room),
	make(chan LobbyReq, 1000),
}
