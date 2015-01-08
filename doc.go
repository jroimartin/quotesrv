// Copyright 2014 The quotesrv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

quotesrv is a "quotes server", that exposes a REST API to add and
list quotes.

Usage:
	quotesrv [flag]

The flags are:
	-addr=":8001": HTTP service address
	-auth=false: enable basic authentication
	-cert="cert.pem": certificate file
	-key="key.pem": private key file
	-pass="s3cr3t": basic auth password
	-quotesfile="quotes.txt": quotes file
	-tls=false: enable TLS
	-user="user": basic auth username

For instance, to run an unauthenticated server over HTTP listening
on IP address 1.1.1.1 and port 8001, you can execute the following
command:

	quotesrv -addr=1.1.1.1:8001

From the client's perspective, the service can be used in the
following way:

Add a new quote (POST request):

	curl http://1.1.1.1:8001/ -d "This is my first quote"

List all quotes (GET request):

	curl http://1.1.1.1:8001/

*/
package main
