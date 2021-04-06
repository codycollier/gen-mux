package main

import (
	"github.com/codycollier/gen-mux/pkg/mux"
)

func main() {
	addr := "127.0.0.1:8888"
	mux.StartMuxServer(addr)
}
