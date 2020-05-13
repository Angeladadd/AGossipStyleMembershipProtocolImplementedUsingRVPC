package simple

const (
	P_FAIL=0.5
	GOSSIP_INTERVAL=1000 //1s发一次
)


type Node struct {
	Membership Membership
	trans chan []Message
	//与传值无关，暂且传个bool
	time chan bool
	timeout chan bool
}


