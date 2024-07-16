package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var port = flag.Int("p", 8080, "port")

func main() {
	// panics are like SEGFAULTs in C. absolutely show-stopping
	// errors that will halt the thread.
	//
	// in order to cleanly log what went wrong, the `defer`
	// statement will run the following function call at the
	// end of the calling function, regardless of whether the
	// thread is panicking.
	defer recoverer()

    flag.Parse()

	// create the same server as before, with a simple `ping` handler
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.HandlerFunc(pingHandler),
	}

	// print something so we know the server is running
	fmt.Printf("listening on port %s\n", server.Addr)

	// the `go` keyword launches a thread. easy-peasy.
	//
	// here we run the server on another thread so we can shut it down
	// from the main thread.
	//
	// likewise, any server panics will be local to the server
	// thread, so in the event of a panic the main thread can still shut
	// down gracefully.
	go launchServer(&server)

	// context.Context provides two primary features, being task cancellation
	// and arbitrary key-value storage. in this case, we use it to "cancel"
	// (shutdown) the server when an interrupt (Ctrl+C) occurs.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// wait for the context to be cancelled.
	<-ctx.Done()

	// create a new context, this time to be cancelled after a 5-second timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shut down the server
	server.Shutdown(ctx)

	fmt.Println("goodbye :)")
}

func launchServer(server *http.Server) {
	// launch the server
	err := server.ListenAndServe()

	// panic on any error (besides the server being shut down).
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

// our ping handler
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong!"))
}

func recoverer() {
	// recover() is a built-in function that returns the value passed to
	// panic(), if the current thread is panicking.
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "unexpected panic occurred: %v", r)
		os.Exit(1)
	}
}
