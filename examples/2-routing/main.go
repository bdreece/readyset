package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type TodoItem struct {
	Content string
	Done    bool
}

var (
	port      = flag.Int("p", 8080, "port")
	todoItems = make([]TodoItem, 0)
)

func main() {
	defer recoverer()
	flag.Parse()

	// these outline the typical REST mapping of create/read/update/delete
	// (CRUD) operations.
	router := http.NewServeMux()
	// returns a list of todo items
	router.HandleFunc("GET /todo-item", getTodoItems)
	// returns a single todo item by id
	router.HandleFunc("GET /todo-item/{id}", getTodoItem)
	// returns status 200 if the item exists, 404 otherwise.
	// this is basically the above GET request, but without sending
	// down the response body. you'll see in the handler code.
	router.HandleFunc("HEAD /todo-item/{id}", getTodoItemExists)
	// creates a new todo item
	router.HandleFunc("POST /todo-item", createTodoItem)
	// replaces or creates a todo item with the given id
	router.HandleFunc("PUT /todo-item/{id}", replaceTodoItem)
	// updates a set of fields on a todo item
	router.HandleFunc("PATCH /todo-item/{id}", updateTodoItem)
	// deletes a todo item
	router.HandleFunc("DELETE /todo-item/{id}", deleteTodoItem)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	fmt.Printf("listening on port %s\n", server.Addr)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	fmt.Println("goodbye :)")
}

// returns the list of todo items.
func getTodoItems(w http.ResponseWriter, r *http.Request) {
	// http.ResponseWriter implements the io.Writer interface,
	// so we can write JSON directly to the response body using
	// a json.Encoder.
	_ = json.NewEncoder(w).Encode(todoItems)
}

// returns a single todo item
func getTodoItem(w http.ResponseWriter, r *http.Request) {
	// get the {id} from /todo-item/{id}, converted to int
	id, err := strconv.Atoi(r.PathValue("id"))
	// return a bad request error if the id is not an int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("`id` must be of type int"))
		return
	}

	// return a not found error if the id is out of range
	if id < 0 || id > len(todoItems)-1 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("`id` not found"))
		return
	}

	item := todoItems[id]
	_ = json.NewEncoder(w).Encode(item)
}

// returns whether a todo item exists
func getTodoItemExists(w http.ResponseWriter, r *http.Request) {
	// get the {id} from /todo-item/{id}, converted to int
	id, err := strconv.Atoi(r.PathValue("id"))
	// return a bad request error if the id is not an int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("`id` must be of type int"))
		return
	}

	// return a not found error if the id is out of range
	if id < 0 || id > len(todoItems)-1 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("`id` not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// creates a todo item
func createTodoItem(w http.ResponseWriter, r *http.Request) {
	// parse form values as strings from request body
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid form"))
		return
	}

	// read form values
	content := r.FormValue("content")
	// convert "done" to bool, send bad request error otherwise
	done, err := strconv.ParseBool(r.FormValue("done"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("`done` must be of type bool"))
		return
	}

	// append new todo item
	todoItems = append(todoItems, TodoItem{
		Content: content,
		Done:    done,
	})

	// create a "Location" header to specify where the new
	// todo item was created (for future GET, HEAD, PUT, etc.).
	location := fmt.Sprintf("%s/todo-item/%d", r.URL.Host, len(todoItems)-1)
	w.Header().Add("Location", location)
	w.WriteHeader(http.StatusCreated)
}

// replaces a todo item
func replaceTodoItem(w http.ResponseWriter, r *http.Request) {
	// get the {id} from /todo-item/{id}, converted to int
	id, err := strconv.Atoi(r.PathValue("id"))
	// return a bad request error if the id is not an int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("`id` must be of type int"))
		return
	}

	// return a not found error if the id is out of range
	if id < 0 || id > len(todoItems)-1 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("`id` not found"))
		return
	}

	// parse form values as strings from request body
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid form"))
		return
	}

	// read form values
	content := r.FormValue("content")
	// convert "done" to bool, send bad request error otherwise
	done, err := strconv.ParseBool(r.FormValue("done"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("`done` must be of type bool"))
		return
	}

	todoItems[id] = TodoItem{
		Content: content,
		Done:    done,
	}

	// create a "Location" header to specify where the new
	// todo item was created (for future GET, HEAD, PUT, etc.).
	location := fmt.Sprintf("%s/todo-item/%d", r.URL.Host, len(todoItems)-1)
	w.Header().Add("Location", location)
	w.WriteHeader(http.StatusOK)
}

// updates a todo item
func updateTodoItem(w http.ResponseWriter, r *http.Request) {
	// get the {id} from /todo-item/{id}, converted to int
	id, err := strconv.Atoi(r.PathValue("id"))
	// return a bad request error if the id is not an int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("`id` must be of type int"))
		return
	}

	// return a not found error if the id is out of range
	if id < 0 || id > len(todoItems)-1 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("`id` not found"))
		return
	}

	// parse form values as strings from request body
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid form"))
		return
	}

	// read form values and update todo item accordingly.

	// update `content`
	content := r.FormValue("content")
	if content != "" {
		todoItems[id].Content = content
	}

	// update `done`
	doneStr := r.FormValue("done")
	if doneStr != "" {
		// convert "done" to bool, send bad request error otherwise
		done, err := strconv.ParseBool(doneStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("`done` must be of type bool"))
			return
		}

		todoItems[id].Done = done
	}

	w.WriteHeader(http.StatusOK)
}

// deletes a todo item
func deleteTodoItem(w http.ResponseWriter, r *http.Request) {
	// get the {id} from /todo-item/{id}, converted to int
	id, err := strconv.Atoi(r.PathValue("id"))
	// return a bad request error if the id is not an int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("`id` must be of type int"))
		return
	}

	// return a not found error if the id is out of range
	if id < 0 || id > len(todoItems)-1 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("`id` not found"))
		return
	}

	// delete the todo item (in Go this is done by appending
	// the two slices of the array that don't contain the id).
	todoItems = append(todoItems[:id], todoItems[id+1:]...)

	w.WriteHeader(http.StatusOK)
}

func recoverer() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "unexpected panic occurred: %v", r)
		os.Exit(1)
	}
}
