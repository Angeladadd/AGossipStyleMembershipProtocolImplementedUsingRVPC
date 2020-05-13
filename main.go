package main

import (
	"time"
	// "fmt"
    "./simple"
    "strconv"
)

func main() {
    nodes := make([]simple.Node, 0)
    for i:=0;i<simple.NODE_NUM;i++ {
        nodes = append(nodes, *simple.NewNode("address"+strconv.Itoa(i)))
    }

    for i:=0;i<simple.NODE_NUM;i++ {
        nodes[i].Others = append(nodes[i].Others, nodes[:i]...)
        nodes[i].Others = append(nodes[i].Others, nodes[i+1:]...)
        // fmt.Println(len(node.Others))
        go nodes[i].Running()
    }

    time.Sleep(time.Second * 10)
}

