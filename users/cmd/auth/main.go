package main

import "users/internal/http"

func main() {
	srv := http.NewServer()
	srv.ConfigServer()
}
