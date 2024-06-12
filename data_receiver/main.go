package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/64bitAryan/go-microservice/types"
	"github.com/gorilla/websocket"
)

func main() {
	receiver := NewDataReceiver()
	http.HandleFunc("/ws", receiver.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgCh chan types.OBUDATA
	conn  *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgCh: make(chan types.OBUDATA, 128),
	}
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU connected")
	for {
		var data types.OBUDATA
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("Received data from OBU [%d] <lat %.2f, lon %2.f>\n", data.OBUID, data.Lat, data.Long)
		dr.msgCh <- data
	}
}
