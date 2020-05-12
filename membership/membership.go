package membership

import (
	"time"
	// "fmt"
)

const (
	BUFSIZE = 1
)

type GossipMessage struct {
	Address string
	Heartbeat int
}

type MembershipCell struct {
	Message GossipMessage
	LocalTime time.Time
	Select chan GossipMessage //\overline{select}
	Update chan GossipMessage //update
	Time chan time.Time //time
}

// MembershipCell的初始化
func GetMembershipCell() *MembershipCell {
	return &MembershipCell{Message:GossipMessage{Address:"hhh"},Select:make(chan GossipMessage), Update:make(chan GossipMessage), Time:make(chan time.Time)}
}

func (cell *MembershipCell) MembershipCellAction() {
	
}