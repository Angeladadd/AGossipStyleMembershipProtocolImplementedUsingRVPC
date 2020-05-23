package simple

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	P_FAIL=0.1
	GOSSIP_INTERVAL=time.Second //1s发一次
	REPAIR_TIME=4*time.Second
	NODE_NUM = 6
	BUFSIZE = 5 //channel buffer size 一般设置为数据中心节点的数目即可
	K = 2
	IS_AUTO = true
)


type Node struct {
	isbad bool
	trans chan []Message
	time chan bool
	timeout chan bool
	accept chan []Message
	deliver chan []Message
	Others []*Node
	//仅打日志及图形化显示使用
	address string
	membership *Membership
}

func NewNode(address string) (instance *Node) {
	instance = new(Node)
	instance.accept = make(chan []Message)
	instance.deliver = make(chan []Message)
	instance.trans = make(chan []Message, BUFSIZE)
	instance.time = make(chan bool)
	instance.timeout = make(chan bool)
	instance.isbad = false
	instance.address = address
	instance.Others = make([]*Node, 0)
	instance.membership = NewMembership(address, instance.accept, instance.deliver)
	fmt.Printf("Initialized Node %p\n", instance)
	return
}

//多路复用实现Nondeterminated Choice
func (node *Node) Fragile() {
	fmt.Printf("Node %p Starts\n", node)
	go node.timer()
	go node.membership.Running()
	node.time <- true
	for {
		if (IS_AUTO) {
			node.Bad()
		}
		if (node.isbad) {
			time.Sleep(REPAIR_TIME)
			continue
		}
		select {
		case message := <- node.trans:
			node.deliver <- message
			for len(node.trans) > 0 {
				message = <- node.trans
				node.deliver <- message
			}
		case <- node.timeout:
			messages := <- node.accept
			node.Gossiping(messages)
		}
	}
}
func (node *Node) Bad() {
	rand.Seed(time.Now().UnixNano())
		r := rand.Float32()
		if r < P_FAIL {
			node.isbad = true
		} else {
			node.isbad = false
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

func (node *Node) Gossiping(messages []Message) {
	rand.Seed(time.Now().UnixNano())
	//随机全排列，取前K个
	perm := rand.Perm(len(node.Others))[:K]
	var targets []*Node
	for _, p := range perm {
		targets = append(targets, node.Others[p])
	}
	str := transmitting(messages, targets)
	fmt.Printf("[SEND] From: %s; To: %s\n", node.address, str)
	node.time <- true
}

func transmitting(messages []Message, targets []*Node) string {
	var str string
	for _, t := range targets {
		if (len(t.trans) == BUFSIZE) {
			continue
		}
		t.trans <- messages
		str+=(t.address+" ")
	}
	return str
}

func (node *Node)ChangeStatus() {
	node.isbad = !node.isbad
}

func (node *Node)Address() string {
	return node.address
}