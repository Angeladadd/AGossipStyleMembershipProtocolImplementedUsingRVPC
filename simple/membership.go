package simple

import (
	"time"
	"fmt"
)

type Message struct {
	Address string
	Heartbeat int
}

type Cell struct {
	Message Message
	LocalTime int64
}

type Membership struct {
	Address string
	Heartbeat int
	MembershipList map[string]Cell
}

func NewMembership(address string) (instance *Membership){
	instance = new(Membership)
	instance.Address = address
	instance.Heartbeat = 0
	instance.MembershipList = make(map[string]Cell,0)
	fmt.Printf("Initialzed Membership %p\n", instance)
	return
}

func (membership *Membership) Deliver(messages []Message) {
	list := membership.MembershipList
	for _, message := range messages {
		if message.Address == membership.Address {
			continue
		}
		if cell, ok := list[message.Address]; !ok || (cell.Message.Heartbeat < message.Heartbeat) {
			//MembershipList中没有这个节点的信息或信息是旧的，增加或更新
			list[message.Address] = Cell{Message:Copy(message), LocalTime:time.Now().UnixNano()}
		} 
	}
	membership.PrintUpdate()
}

func (membership *Membership) Accept() (messages []Message) {
	list := membership.MembershipList
	messages = make([]Message, 0)
	membership.Heartbeat++
	messages = append(messages, Message{membership.Address, membership.Heartbeat})
	for _, cell := range list {
		messages = append(messages, Copy(cell.Message))
	}
	return
}