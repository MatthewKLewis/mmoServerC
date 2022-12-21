package main

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	socket websocket.Conn
	uuid   string
	name   string
	x      float64
	y      float64
	z      float64
	rotY   float64
}

func (p *Player) generateMovementPacketFromState() Packet {
	return Packet{
		header: "MOVEMENT",
		uuid:   p.uuid,
		x:      p.x,
		y:      p.y,
		z:      p.z,
		rotY:   p.rotY,
	}
}
