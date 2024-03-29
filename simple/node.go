package simple

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	P_FAIL=0.1 //自动失效概率
	GOSSIP_INTERVAL=time.Second //Gossip周期
	REPAIR_TIME=4*time.Second //失效节点修复时间
	NODE_NUM = 5 //节点数目
	BUFSIZE = 4 //channel buffer size 一般设置为数据中心节点的数目即可
	K = 2 //每次Gossip的目标节点数
	IS_AUTO = true //节点是否自动随机失效
)

type Nil struct {}

type Node struct {
	isbad bool
	trans chan []Message
	time chan Nil
	timeout chan Nil
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
	instance.time = make(chan Nil)
	instance.timeout = make(chan Nil)
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
	node.time <- Nil{}
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
			rand.Seed(time.Now().UnixNano())
			perm := rand.Perm(len(node.Others))[:K]
			var targets []*Node
			for _, p := range perm {
				targets = append(targets, node.Others[p])
			}
			node.Gossiping(messages, targets)
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
			node.timeout <- Nil{}
		}
	}
}

func (node *Node) Gossiping(messages []Message, targets []*Node) {
	var str string
	for _, t := range targets {
		if (len(t.trans) == BUFSIZE) {
			continue
		}
		t.trans <- messages
		str+=(t.address+" ")
	}
	fmt.Printf("[SEND] From: %s; To: %s\n", node.address, str)
	node.time <- Nil{}
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
