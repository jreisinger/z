// Guessdomain tries to enumerate subdomains of a given domain.
//
//	go run guessdomain.go go.dev < wordlist.txt
//
// For more wordlists see
// https://github.com/danielmiessler/SecLists/tree/master/Discovery/DNS
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jreisinger/z"
	"github.com/miekg/dns"
)

const dnsServer = "8.8.8.8:53" // DNS server address to use for lookups

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "guessdomain: missing domain to perform guessing against, e.g. example.com")
		os.Exit(1)
	}

	g := &guess{
		domain:     os.Args[1],
		dnsSrvAddr: dnsServer,
	}
	z.Run(g)
}

type guess struct {
	domain     string
	dnsSrvAddr string
}

func (g *guess) Make(line string) z.Task {
	subdomain := line
	t := &domain{
		fqdn:          fmt.Sprintf("%s.%s", subdomain, g.domain),
		dnsServerAddr: g.dnsSrvAddr,
	}
	return t
}

type domain struct {
	dnsServerAddr string
	fqdn          string
	ipAddrs       []string
}

func (d *domain) Process() {
	results := lookup(d.fqdn, d.dnsServerAddr)
	d.ipAddrs = results
}

func (d *domain) Print() {
	for _, ip := range d.ipAddrs {
		fmt.Printf("%s\t%s\n", d.fqdn, ip)
	}
}

func lookup(fqdn, serverAddr string) (ipAddrs []string) {
	var cfqdn = fqdn // Don't modify the original.
	for {
		cnames, err := lookupCNAME(cfqdn, serverAddr)
		if err == nil && len(cnames) > 0 {
			cfqdn = cnames[0]
			continue // We have to process the next CNAME.
		}
		ips, err := lookupA(cfqdn, serverAddr)
		if err != nil {
			break // There are no A records for this hostname.
		}
		ipAddrs = append(ipAddrs, ips...)
		break // We have processed all the results.
	}
	return
}

func lookupA(fqdn, serverAddr string) ([]string, error) {
	var m dns.Msg
	var ips []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	r, err := dns.Exchange(&m, serverAddr)
	if err != nil {
		return ips, err
	}
	if len(r.Answer) < 1 {
		return ips, errors.New("no answer")
	}
	for _, answer := range r.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

func lookupCNAME(fqdn, serverAddr string) ([]string, error) {
	var m dns.Msg
	var fqdns []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, serverAddr)
	if err != nil {
		return fqdns, err
	}
	if len(in.Answer) < 1 {
		return fqdns, errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, c.Target)
		}
	}
	return fqdns, nil
}
