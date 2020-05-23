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
	address string
	heartbeat int
	membershipList map[string]Cell
	accept chan []Message
	deliver chan []Message
}

func NewMembership(address string, accept chan []Message, deliver chan[]Message) (instance *Membership){
	instance = new(Membership)
	instance.address = address
	instance.heartbeat = 0
	instance.membershipList = make(map[string]Cell,0)
	instance.accept = accept
	instance.deliver = deliver
	fmt.Printf("Initialzed Membership %p\n", instance)
	return
}

func (membership *Membership) Running() {
	for {
		select {
		case messages := <- membership.deliver:
			membership.Deliver(messages)
		case membership.accept <- membership.Accept():

		}
	}
}

func (membership *Membership) Deliver(messages []Message) {
	list := membership.membershipList
	for _, message := range messages {
		if message.Address == membership.address {
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
	list := membership.membershipList
	messages = make([]Message, 0)
	membership.heartbeat++
	messages = append(messages, Message{membership.address, membership.heartbeat})
	for _, cell := range list {
		messages = append(messages, Copy(cell.Message))
	}
	return
}