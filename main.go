package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/jroimartin/orujo"
	"github.com/jroimartin/orujo-handlers/basic"
	olog "github.com/jroimartin/orujo-handlers/log"
)

var (
	addr       = flag.String("addr", ":8001", "HTTP service address")
	quotesFile = flag.String("quotesfile", "quotes.txt", "file where quotes will be stored")
	auth       = flag.Bool("auth", false, "enable basic authentication")
	user       = flag.String("user", "user", "username used for basic authentication")
	pass       = flag.String("pass", "s3cr3t", "password used for basic authentication")
	tls        = flag.Bool("tls", false, "enable TLS")
	certFile   = flag.String("cert", "cert.pem", "certificate file")
	keyFile    = flag.String("key", "key.pem", "private key file")
	re         = regexp.MustCompile(`[\r\n]+`)
)

var mutex sync.RWMutex

func main() {
	flag.Parse()

	s := orujo.NewServer(*addr)

	logger := log.New(os.Stdout, "[quotesrv] ", log.LstdFlags)
	logHandler := olog.NewLogHandler(logger, logLine)

	var basicHandler http.Handler
	if *auth {
		basicHandler = basic.NewBasicHandler("Quotes System", *user, *pass)
	} else {
		basicHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}

	s.RouteDefault(http.NotFoundHandler(), orujo.M(logHandler))

	s.Route(`^/$`,
		basicHandler,
		http.HandlerFunc(listQuotes),
		orujo.M(logHandler),
	).Methods("GET")

	s.Route(`^/$`,
		basicHandler,
		http.HandlerFunc(addQuote),
		orujo.M(logHandler),
	).Methods("POST")

	if *tls {
		logger.Fatalln(s.ListenAndServeTLS(*certFile, *keyFile))
	} else {
		logger.Fatalln(s.ListenAndServe())
	}
}

func listQuotes(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	f, err := os.Open(*quotesFile)
	if err != nil {
		errorResponse(w, err)
		return
	}
	defer f.Close()

	if _, err = io.Copy(w, f); err != nil {
		errorResponse(w, err)
		return
	}
}

func addQuote(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	f, err := os.OpenFile(*quotesFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		errorResponse(w, err)
		return
	}
	defer f.Close()

	bQuote, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}
	quote := re.ReplaceAllString(string(bQuote), " ") + "\n"

	if _, err = f.WriteString(quote); err != nil {
		errorResponse(w, err)
		return
	}
}

func errorResponse(w http.ResponseWriter, err error) {
	orujo.RegisterError(w, err)
	w.WriteHeader(http.StatusInternalServerError)
}

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`
