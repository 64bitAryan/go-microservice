package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/64bitAryan/go-microservice/types"
	"github.com/gorilla/websocket"
)

func main() {

	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", receiver.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgCh chan types.OBUDATA
	conn  *websocket.Conn
	prod  DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)
	p, err = NewkafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMidleware(p)
	return &DataReceiver{
		msgCh: make(chan types.OBUDATA, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUDATA) error {
	return dr.prod.ProduceData(data)
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
		if err := dr.produceData(data); err != nil {
			log.Println(err)
			continue
		}
	}
}
