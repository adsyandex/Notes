// main.go
package main

import (
    "log"
    "net/http"
)

func main() {
    db := InitDB()
    InitHandlers(db)

    log.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
