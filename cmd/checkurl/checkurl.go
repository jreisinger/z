// Checkurl finds out whether URLs return OK (200) status.
package main

import (
	"fmt"
	"net/http"

	"github.com/jreisinger/z"
)

type urls struct{}

func (*urls) Make(line string) z.Task {
	return &resource{url: line}
}

type resource struct {
	url    string
	status bool
}

func (r *resource) Process() {
	resp, err := http.Get(r.url)
	if err != nil {
		return
	}
	if resp.StatusCode == http.StatusOK {
		r.status = true
	}
}

func (r *resource) Print() {
	status := map[bool]string{
		true:  "OK",
		false: "NOTOK",
	}
	fmt.Printf("%-5s %s\n", status[r.status], r.url)
}

func main() {
	z.Run(&urls{})
}
