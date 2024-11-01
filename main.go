package main

import (
	"log"
	"net/http"
	"os"
)



func main() {
    args := os.Args[1:]
    if len(args) == 0 {
        log.Fatal("Missing branch name")
    }
    branchName := args[0]
    log.Printf("Creating public hostname for branch %s", branchName)

    os.

    resp, err := http.MethodPut()

}
