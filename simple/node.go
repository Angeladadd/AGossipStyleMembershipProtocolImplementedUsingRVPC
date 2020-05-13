package simple

import (
	"math/rand"
	"time"
)

const (
	P_FAIL=0.5
	GOSSIP_INTERVAL=1000 //1s发一次
	BUFSIZE = 4 //channel buffer size 一般设置为数据中心节点的数目即可
	K = 2
)


type Node struct {
	IsBad bool
	Membership *Membership
	//论文中的多个trans只是表示了网络了链接的路径，实际只用一个trans来接收消息就可以了，\overline{trans}不需要显式的表示
	Trans chan []Message
	//与传值无关，暂且传个bool
	Time chan bool
	Timeout chan bool
}

func NewNode(address string) (instance *Node) {
	instance = new(Node)
	instance.Membership = NewMembership(address)
	instance.Trans = make(chan []Message, BUFSIZE)
	instance.Time = make(chan bool)
	instance.Timeout = make(chan bool)
	instance.IsBad = false
	return
}

func (node *Node) Deliver(messages []Message) {
	node.Membership.Deliver(messages)
}

func (node *Node) Gossiping(others []Node) {
	messages := node.Membership.Accept()
	cpy := make([]Node, len(others))
	copy(cpy, others)
	//用rand来处理全排列
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cpy), func(i int, j int){cpy[i], cpy[j] = cpy[j], cpy[i]})
	transmitting(messages, others[:K])
}

func transmitting(messages []Message, perms []Node) {
	for _, p := range perms {
		p.Trans <- messages
	}
}

