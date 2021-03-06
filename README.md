# quotesrv

Quotes server

## Installation

`go get github.com/jroimartin/quotesrv`

## Basic usage

### Server

Command's help:

```
$ quotesrv -h
Usage of quotesrv:
  -addr=":8001": HTTP service address
  -auth=false: enable basic authentication
  -cert="cert.pem": certificate file
  -key="key.pem": private key file
  -pass="s3cr3t": basic auth password
  -quotesfile="quotes.txt": quotes file
  -tls=false: enable TLS
  -user="user": basic auth username
```

Run an unauthenticated server over HTTP listening on IP address 1.1.1.1 and
port 8001:

`$ quotesrv -addr=1.1.1.1:8001`

### Client

Add a new quote:

`$ curl http://1.1.1.1:8001/ -d "This is my first quote"`

List all quotes:

`$ curl http://1.1.1.1:8001/`
