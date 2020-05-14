package simple

import (
	"fmt"
	"bytes"
)

func Copy(m Message) Message {
	return Message{m.Address, m.Heartbeat}
}

func (c Cell) String() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{Address:%s,Heartbeat:%d,LocalTime:%d}", c.Message.Address, c.Message.Heartbeat, c.LocalTime.UnixNano())
	return b.String()
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