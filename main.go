package main

import (
	// "time"
	"fmt"
    "./simple"
    "strconv"
    "log"
    "net/http"
    "bytes"
    "io/ioutil"
)

var nodes []*simple.Node

func main() {
    nodes = make([]*simple.Node, 0)
    for i:=0;i<simple.NODE_NUM;i++ {
        nodes = append(nodes, simple.NewNode("address"+strconv.Itoa(i)))
    }

    for i:=0;i<simple.NODE_NUM;i++ {
        nodes[i].Others = append(nodes[i].Others, nodes[:i]...)
        nodes[i].Others = append(nodes[i].Others, nodes[i+1:]...)
        // fmt.Println(len(node.Others))
        go nodes[i].Running()
    }

    for {
        http.HandleFunc("/index", handler)
        http.HandleFunc("/membership", membership)
        log.Fatal(http.ListenAndServe("localhost:8080", nil))
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("templates/index.html")
    fmt.Fprint(w, string(body))
}

func membership(w http.ResponseWriter, r *http.Request) {
    b := new(bytes.Buffer)
    for _, node := range nodes {
        b.WriteString(node.Membership.PrintUpdate()+"\n")
    } 
    fmt.Fprint(w, b.String())
}

