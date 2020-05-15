package tools

import (
	"fmt"
    "time"
    "bytes"
	// "../simple"
)

// func main() {
//     ProgressBar(100 * time.Millisecond)
// }

func ProgressBar(delay time.Duration) {
    var buffer bytes.Buffer
    var cnt int
    for ;cnt<100; {
        fmt.Printf("\r[%s>%d%%]",buffer.String(),cnt)
        buffer.WriteString("=")
        cnt++
        time.Sleep(delay)
    }
    fmt.Printf("\r[%s>%d%%]\n",buffer.String(), cnt)
}

func Spinner(delay time.Duration) {
    for {
        for _, r := range `-\|/` {
            fmt.Printf("\r%c", r)
            time.Sleep(delay)
        }
    }
}

