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
	fmt.Fprintf(b, "{Address:%s,Heartbeat:%d,LocalTime:%d}", c.Message.Address, c.Message.Heartbeat, c)
	return b.String()
}

type Info struct {
	Membership Membership
	LocalTime int64
	IsBad bool
}

func (node Node) Info() Info {
	return Info{Membership:*node.Membership, LocalTime:time.Now().UnixNano(), IsBad:node.IsBad}
}

func (membership *Membership) PrintUpdate() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "[UPDATE] {Address:%s, MembershipList:[",membership.Address)
	for _, value := range membership.MembershipList {
        fmt.Fprintf(b, "%s,", value.String())
	}
	b.WriteString("\b]}\n")
	// fmt.Printf(b.String())
	return b.String()
}

