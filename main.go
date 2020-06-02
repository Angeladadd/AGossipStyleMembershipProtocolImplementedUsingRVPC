package main

import (
	// "encoding/json"
	// "fmt"
    "./simple"
    "strconv"
    // "log"
    // "net/http"
    // "io/ioutil"
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
        go nodes[i].Fragile()
    }

    for {
        // http.HandleFunc("/index", index)
        // http.HandleFunc("/membership", membership)
        // http.HandleFunc("/change_status", changeStatus)
        // log.Fatal(http.ListenAndServe("localhost:8080", nil))
    }
}

// func index(w http.ResponseWriter, r *http.Request) {
//     body, _ := ioutil.ReadFile("templates/index.html")
//     fmt.Fprint(w, string(body))
// }

// func membership(w http.ResponseWriter, r *http.Request) {
//     header := w.Header()
//     header.Add("Content-Type","application/json")
//     w.WriteHeader(http.StatusOK)
//     var s []simple.Info
//     for _, node := range nodes {
//         s = append(s, node.Info())
//     }
//     b, err := json.Marshal(s)
//     if err!= nil {
//         log.Println("marshal error")
//     }
//     fmt.Fprintf(w, string(b))
// }

// func changeStatus(w http.ResponseWriter, r *http.Request) {
//     vars := r.URL.Query() 
//     address := vars["address"][0]
//     for _, node := range nodes {
//         if (node.Address() == address) {
//             node.ChangeStatus()
//         }
//     }
//     header := w.Header()
//     header.Add("Content-Type","application/json")
//     w.WriteHeader(http.StatusOK)
// }

