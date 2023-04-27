// Lookupip looks up IP addresses of hosts using the local resolver.
package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/jreisinger/z"
)

type hosts struct{}

func (hosts) Make(line string) z.Task {
	return &host{name: line}
}

type host struct {
	name string
	err  error
	ips  []net.IP
}

func (h *host) Print() {
	var ips []string
	for _, ip := range h.ips {
		ips = append(ips, ip.String())
	}
	if h.err != nil {
		fmt.Fprintf(os.Stderr, "z: %v", h.err)
	}
	fmt.Printf("%-15s %s\n", h.name, strings.Join(ips, ", "))
}

func (h *host) Process() {
	ips, err := net.LookupIP(h.name)
	if err != nil {
		h.err = err
		return
	}
	h.ips = ips
}

func main() {
	z.Run(&hosts{})
}
