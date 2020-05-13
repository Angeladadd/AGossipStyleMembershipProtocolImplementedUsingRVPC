package main

import (
	"./simple"
)

func main() {
    var membership = simple.NewMembership("hhh")
    var messages = make([]simple.Message, 0)
    messages = append(messages, simple.Message{"1",6})
    messages = append(messages, simple.Message{"2",3})
    membership.Deliver(messages)
    messages[0].Heartbeat = 3
    membership.Deliver(messages)
    messages[1].Heartbeat = 7
    membership.Deliver(messages)
}

