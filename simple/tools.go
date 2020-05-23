package simple

import (
	"fmt"
	"bytes"
	"time"
)

func Copy(m Message) Message {
	return Message{m.Address, m.Heartbeat}
}

func (c Cell) String() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{Address:%s,Heartbeat:%d,LocalTime:%d}", c.Message.Address, c.Message.Heartbeat, c.LocalTime)
	return b.String()
}

type Info struct {
	Address string
	Heartbeat int
	MembershipList map[string]Cell
	LocalTime int64
	IsBad bool
}

func (node Node) Info() Info {
	return Info{Address:node.address, Heartbeat:node.membership.heartbeat, MembershipList:node.membership.membershipList, LocalTime:time.Now().UnixNano(), IsBad:node.isbad}
}

func (membership *Membership) PrintUpdate() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "[UPDATE] {Address:%s, MembershipList:[",membership.address)
	for _, value := range membership.membershipList {
        fmt.Fprintf(b, "%s,", value.String())
	}
	b.WriteString("\b]}\n")
	fmt.Printf(b.String())
	return b.String()
}

