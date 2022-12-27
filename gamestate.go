package main

import (
	"fmt"
	"time"
)

type Gamestate struct {
	currentTime time.Time
	players     []Player
}

// STATE MODIFYING
func (gS *Gamestate) updatePlayerPositions(p Packet) {
	var found = false
	for i := 0; i < len(gS.players); i++ {
		if gS.players[i].uuid == p.uuid {
			found = true
			gS.players[i].x = p.x
			gS.players[i].y = p.y
			gS.players[i].z = p.z
			gS.players[i].rotY = p.rotY
		}
	}

	if !found {
		fmt.Println("Adding to List")
		gS.players = append(gS.players, Player{
			uuid:   p.uuid,
			socket: p.source,
			x:      p.x,
			y:      p.y,
			z:      p.z,
			rotY:   p.rotY,
		})
	}
}

func (gS *Gamestate) resolveAttack(p Packet) {
	fmt.Println("Attack")
}

func (gS *Gamestate) removePlayerFromList(p Packet) {
	fmt.Println("Removing from List")
	var index = -1
	for i := 0; i < len(gS.players); i++ {
		if gS.players[i].uuid == p.uuid {
			index = i
		}
	}
	if index != -1 {
		gS.players = remove(gS.players, index)
	}
}

// OUTGOING
func (gS *Gamestate) sendPlayerPositions(time time.Time) {
	for i := 0; i < len(gS.players); i++ {
		for j := 0; j < len(gS.players); j++ {
			if gS.players[i].uuid != gS.players[j].uuid {
				outPacket := gS.players[j].generateMovementPacketFromState()
				gS.players[i].socket.WriteMessage(1, outPacket.serialize())
			}
		}
	}
}
func (gS *Gamestate) sendPlayerChat(p Packet) {
	for i := 0; i < len(gS.players); i++ {
		for j := 0; j < len(gS.players); j++ {
			if gS.players[i].uuid != gS.players[j].uuid {
				gS.players[i].socket.WriteMessage(1, p.serialize())
			}
		}
	}
}
func (gS *Gamestate) sendPlayerDisconnect(p Packet) {
	fmt.Println("Sending Disconnect Packet to Remaining Players")
	for i := 0; i < len(gS.players); i++ {
		gS.players[i].socket.WriteMessage(1, p.serialize())
	}
}
