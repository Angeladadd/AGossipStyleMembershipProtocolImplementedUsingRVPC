package simple

import (
	"time"
)

const (
	BUFSIZE = 1
)

type Message struct {
	Address string
	Heartbeat int
}

/* 对于大型数据中心，我们可能会有专门的程序来处理Membership
 * 对每个MembershipCell可能会单独开线程处理读与写（这里会涉及读写锁，因为整个协议中读多写少）
 * 由于只有一台落后的4核机器，这里做串行化处理了
 */
type Cell struct {
	Message Message
	LocalTime time.Time
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
	instance.MembershipList = make(map[string]Cell)
	return
}

func Copy(m Message) Message {
	return Message{m.Address, m.Heartbeat}
}

func (membership *Membership) Deliver(messages []Message) {
	list := membership.MembershipList
	for _, message := range messages {
		if message.Address == membership.Address {
			continue
		}
		if cell, ok := list[message.Address]; !ok || (cell.Message.Heartbeat < message.Heartbeat) {
			//MembershipList中没有这个节点的信息或信息是旧的，增加或更新
			list[message.Address] = Cell{Message:Copy(message), LocalTime:time.Now()}
		} 
	}
}

func (membership *Membership) Accept() (messages []Message) {
	list := membership.MembershipList
	messages = make([]Message, len(list)+1)
	membership.Heartbeat++
	messages = append(messages, Message{membership.Address, membership.Heartbeat})
	for _, cell := range list {
		messages = append(messages, Copy(cell.Message))
	}
	return
}