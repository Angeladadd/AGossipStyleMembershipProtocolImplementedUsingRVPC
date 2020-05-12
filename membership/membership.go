package membership

import (
	"time"
)

const (
	BUFSIZE = 1
)

type GossipMessage struct {
	Address string
	Heartbeat int
}

/* 对于大型数据中心，我们可能会有专门的程序来处理Membership
 * 对每个MembershipCell可能会单独开线程处理读与写（这里会涉及读写锁，因为整个协议中读多写少）
 * 由于只有一台落后的4核机器，这里做串行化处理了
 */
type MembershipCell struct {
	Message GossipMessage
	LocalTime time.Time
}

type Membership struct {
	Address string
	Heartbeat int
	MembershipList map[string]MembershipCell
}

func NewMembership(address string) (instance *Membership){
	instance = new(Membership)
	instance.Address = address
	instance.Heartbeat = 0
	instance.MembershipList = make(map[string]MembershipCell)
	return
}

func Copy(m GossipMessage) GossipMessage {
	return GossipMessage{m.Address, m.Heartbeat}
}

func (membership *Membership)Deliver(messages []GossipMessage) {
	list := membership.MembershipList
	for _, message := range messages {
		if message.Address == membership.Address {
			continue
		}
		if cell, ok := list[message.Address]; !ok || (cell.Message.Heartbeat < message.Heartbeat) {
			//MembershipList中没有这个节点的信息或信息是旧的，增加或更新
			list[message.Address] = MembershipCell{Message:Copy(message), LocalTime:time.Now()}
		} 
	}
}

func (membership *Membership) Accept() (messages []GossipMessage) {
	list := membership.MembershipList
	messages = make([]GossipMessage, len(list)+1)
	membership.Heartbeat++
	messages = append(messages, GossipMessage{membership.Address, membership.Heartbeat})
	for _, cell := range list {
		messages = append(messages, Copy(cell.Message))
	}
	return
}