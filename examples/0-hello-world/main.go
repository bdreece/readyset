package main

import (
    // flag provides command-line argument parsing
    //  e.g. ./hello-world -p 8080
	"flag"
    // fmt provides formatting and printing utilities
	"fmt"
    // net/http provides http functionality
	"net/http"
)

// flag.Int(...) declares a command line argument
// to parse at startup. In this case, it is
// for the port our server will listen on.
var port = flag.Int("p", 8080, "port")

func main() {
    // parse the command line arguments
    flag.Parse()

    // http.ServeMux is a type of catch-all
    // handler. its main job is to forward
    // requests (multiplex, or mux) to the
    // appropriate sub-handler based on the
    // method and path.
	router := http.NewServeMux()

    // register `helloWorld` to run on requests
    // with method "GET" and path "/".
	router.HandleFunc("GET /", helloWorld)

    // http.Server provides the TCP socket listener
    // that will accept all incoming HTTP connections
    // on the provided address.
    server := http.Server{
        // e.g. :8080
        Addr: fmt.Sprintf(":%d", *port),
        // a server uses one handler, hence the
        // http.ServeMux router above needed to
        // route to multiple handlers
        Handler: router,
    }


    // print something so we know the server is running
    fmt.Printf("server listening on %s\n", server.Addr)

    // finally, listen on the specified port. this
    // will block until either a fatal error occurs
    // or server.Shutdown(...) is called.
    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }

    // leave little goodies behind. unfortunately this
    // will never print, see [1-shutdown].
    fmt.Println("goodbye :)")
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
    _, _ = w.Write([]byte("hello, world!"))
}
