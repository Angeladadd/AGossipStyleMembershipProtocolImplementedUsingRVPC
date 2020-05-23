package simple

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	P_FAIL=0.1
	GOSSIP_INTERVAL=time.Second //1s发一次
	REPAIR_TIME=2*time.Second
	NODE_NUM = 5
	BUFSIZE = 4 //channel buffer size 一般设置为数据中心节点的数目即可
	K = 2
)


type Node struct {
	isbad bool
	membership *Membership
	//论文中的多个trans只是表示了网络了链接的路径，实际只用一个trans来接收消息就可以了，\overline{trans}不需要显式的表示
	trans chan []Message
	//与传值无关，暂且传个bool
	time chan bool
	timeout chan bool
	Others []*Node
}

func NewNode(address string) (instance *Node) {
	instance = new(Node)
	instance.membership = NewMembership(address)
	instance.trans = make(chan []Message, BUFSIZE)
	instance.time = make(chan bool)
	instance.timeout = make(chan bool)
	instance.isbad = false
	instance.Others = make([]*Node, 0)
	fmt.Printf("Initialized Node %p\n", instance)
	return
}

//多路复用实现Nondeterminated Choice
func (node *Node) Fragile() {
	fmt.Printf("Node %p Starts\n", node)
	go node.timer()
	node.time <- true
	for {
		rand.Seed(time.Now().UnixNano())
		r := rand.Float32()
		if r < P_FAIL {
			node.isbad = true
		} else {
			node.isbad = false
		}
		if (node.isbad) {
			time.Sleep(REPAIR_TIME)
			continue
		}
		select {
		case message := <- node.trans:
			node.Deliver(message)
			for len(node.trans) > 0 {
				message = <- node.trans
				node.Deliver(message)
			}
		case <- node.timeout:
			messages := node.membership.Accept()
			node.Gossiping(messages)
		}
	}
}

func (node *Node) timer() {
	for {
		select {
		case <- node.time:
			time.Sleep(GOSSIP_INTERVAL)
			node.timeout <- true
		}
	}
}


func (node *Node) Deliver(messages []Message) {
	node.membership.Deliver(messages)
}

func (node *Node) Gossiping(messages []Message) {
	rand.Seed(time.Now().UnixNano())
	//随机全排列，取前K个
	perm := rand.Perm(len(node.Others))[:K]
	var targets []*Node
	for _, p := range perm {
		targets = append(targets, node.Others[p])
	}
	str := transmitting(messages, targets)
	fmt.Printf("[SEND] From: %s; To: %s\n", node.membership.Address, str)
	node.time <- true
}

func transmitting(messages []Message, targets []*Node) string {
	var str string
	for _, t := range targets {
		t.trans <- messages
		str+=(t.membership.Address+" ")
	}
	return str
}