Z is a simple Go package that allows you easily build CLI tools that concurrently process lines from STDIN.

See `cmd` folder for such tools.

To use the tools, for example:

```
$ go doc cmd/lookupip/lookupip.go
Lookupip looks up IP addresses of hosts using the local resolver.
$ go run cmd/lookupip/lookupip.go < /tmp/hosts
golang.org      142.251.209.17, 2a00:1450:4002:410::2011
perl.com        151.101.2.132, 151.101.66.132, 151.101.194.132, 151.101.130.132
python.org      151.101.65.168, 151.101.193.168, 151.101.129.168, 151.101.1.168
```

or 

```
$ go install cmd/lookupip/lookupip.go # installs into ~/go/bin by default
$ lookupip < /tmp/hosts
```

or

```
$ GOOS=linux GOARCH=arm64 go build cmd/lookupip/lookupip.go # go tool dist list
$ scp ./lookupip user@raspberry.net:
$ ssh user@raspberry.net
raspberry$ ./lookupip < /tmp/hosts
```

Stolen from [jgc](https://github.com/cloudflare/jgc-talks/tree/master/dotGo/2014) ([video](https://youtu.be/woCg2zaIVzQ)).
