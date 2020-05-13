package simple

import (
	"math/rand"
	"time"
)

const (
	P_FAIL=0.5
	GOSSIP_INTERVAL=1000 //1s发一次
	NODE_NUM = 5
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
	Others []Node
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

func (node *Node) Bad() {
	rand.Seed(time.Now().UnixNano())
	r := rand.Int31n(2) // 0,1
	if r < 1 {
		node.IsBad = true
	} else {
		node.IsBad = false
	}
}

func (node *Node) Deliver(messages []Message) {
	node.Membership.Deliver(messages)
}

func (node *Node) Gossiping() {
	messages := node.Membership.Accept()
	rand.Seed(time.Now().UnixNano())
	//随机全排列，取前K个
	perm := rand.Perm(len(node.Others))[:K]
	var targets []Node
	for _, p := range perm {
		targets = append(targets, node.Others[p])
	}
	transmitting(messages, targets)
}

func transmitting(messages []Message, targets []Node) {
	for _, t := range targets {
		t.Trans <- messages
	}
}