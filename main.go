package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func lazyApproveOrigin(r *http.Request) bool { return true }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     lazyApproveOrigin,
}

var incoming = make(chan Packet)
var ticker = time.NewTicker(100 * time.Millisecond)
var gamestate = Gamestate{}

const host = "0.0.0.0"
const port = "8000"

func main() {
	fmt.Println("Starting Server")
	go gameLoop()
	go http.HandleFunc("/", socketHandler)
	go log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func gameLoop() {
	//var timeLastSent = time.Now()
	for {
		select {
		case iPacket := <-incoming:
			switch iPacket.header {
			case "MOVEMENT":
				gamestate.updatePlayerPositions(iPacket)
			case "ATTACK":
				gamestate.resolveAttack(iPacket)
			case "DISCONNECT":
				gamestate.removePlayerFromList(iPacket)
				gamestate.sendPlayerDisconnect(iPacket) //a sender, not a state modifier, but no racecons...
			case "CHAT":
				gamestate.sendPlayerChat(iPacket) //another sender.
			default:
				fmt.Println("Couldn't Switchboard Packet!", iPacket)
			}
		case tic := <-ticker.C:
			//fmt.Println(time.Since(timeLastSent))
			//timeLastSent = time.Now()
			gamestate.sendPlayerPositions(tic)
		}
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade Connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//Get Player and User Information from API (SIMULATED)
	time.Sleep(50 * time.Millisecond)

	fmt.Println("New Connection")
	var uuid = uuid.New().String()

	err = conn.WriteMessage(1, []byte("SPAWN|"+uuid+"|Matthew|0,5,0,0"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// PROCESS MESSAGES
	for {
		var retPack = Packet{}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			//HANDLE INVALID MESSAGES WITH DISCONNECTION PACKET
			fmt.Println("Client Disconnect")
			retPack.uuid = uuid
			retPack.header = "DISCONNECT"
			incoming <- retPack
			return
		}

		//SWITCHBOARD
		var msgArray = strings.Split(string(msg), "|")
		var coordinateArray = strings.Split(msgArray[2], ",")
		switch msgArray[0] {
		case "MOVEMENT":
			retPack.uuid = uuid
			retPack.header = "MOVEMENT"
			retPack.source = *conn
			retPack.x, _ = strconv.ParseFloat(coordinateArray[0], 64)
			retPack.y, _ = strconv.ParseFloat(coordinateArray[1], 64)
			retPack.z, _ = strconv.ParseFloat(coordinateArray[2], 64)
			retPack.rotY, _ = strconv.ParseFloat(coordinateArray[3], 64)
			break
		case "ATTACK":
			retPack.uuid = uuid
			retPack.header = "ATTACK"
			retPack.source = *conn
			retPack.x, _ = strconv.ParseFloat(coordinateArray[0], 64)
			retPack.y, _ = strconv.ParseFloat(coordinateArray[1], 64)
			retPack.z, _ = strconv.ParseFloat(coordinateArray[2], 64)
			retPack.rotY, _ = strconv.ParseFloat(coordinateArray[3], 64)
			break
		case "CHAT":
			retPack.uuid = uuid
			retPack.header = "CHAT"
			retPack.source = *conn
			retPack.x, _ = strconv.ParseFloat(coordinateArray[0], 64)
			retPack.y, _ = strconv.ParseFloat(coordinateArray[1], 64)
			retPack.z, _ = strconv.ParseFloat(coordinateArray[2], 64)
			retPack.other = "Hello World!"
			break
		}

		//SEND
		incoming <- retPack
	}
}
