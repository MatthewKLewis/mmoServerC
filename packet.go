package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Packet struct {
	header string
	uuid   string
	source websocket.Conn
	x      float64
	y      float64
	z      float64
	rotY   float64
	other  string
}

func (p *Packet) serialize() []byte {
	return []byte(p.header + "|" + fmt.Sprint(p.uuid) + "|" + "Matthew!" + "|" + fmt.Sprint(p.x) + "," + fmt.Sprint(p.y) + "," + fmt.Sprint(p.z) + "," + fmt.Sprint(p.rotY))
}
