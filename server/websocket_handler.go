package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fsufitch/slickrlink-transfer-node/protobufs"
	"github.com/fsufitch/slickrlink-transfer-node/uploadsession"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var defaultUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type clientWebsocketHandler struct {
	upgrader websocket.Upgrader
}

func writeClientMessages(conn *websocket.Conn, s *uploadsession.TransferSession) {
	for msg := range s.OutgoingMessages {
		msg.Timestamp = time.Now().Unix()
		data, err := proto.Marshal(msg)
		if err != nil {
			log.Printf("Error marshaling message %s", err.Error())
			continue
		}

		err = conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			s.FatalErrors <- err
			return
		}
	}
}

func receiveClientMessages(conn *websocket.Conn, s *uploadsession.TransferSession) {
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			s.FatalErrors <- err
			return
		}

		if msgType != websocket.BinaryMessage {
			s.FatalErrors <- errors.New("Unexpected non-binary message")
			return
		}

		msg := &protobufs.ClientToTransferNodeMessage{}
		err = proto.Unmarshal(p, msg)
		if err != nil {
			s.FatalErrors <- err
			return
		}

		s.IncomingMessages <- msg
	}
}

func (h clientWebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Println(err)
		return
	}

	transferSession := uploadsession.NewTransferSession()
	defer transferSession.Destroy()

	go receiveClientMessages(conn, transferSession)
	go writeClientMessages(conn, transferSession)
	go transferSession.HandleMessages()

connectionLoop:
	for {
		select {
		case err = <-transferSession.FatalErrors:
			log.Printf("Fatal error in websocket connection %s", err.Error())
			break connectionLoop
		}
	}

}
