package main

import (
	"fmt"
	"internal/connection"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Message struct {
	Message string `json:"message"`
}

func main() {
	err := connection.InitDB()
	if err != nil {
		log.Printf("failed to init DB: %v", err)
	}
	defer connection.DB.Close()

	router := httprouter.New()
	router.POST("/lamp/get/:id", getLampHandler)
	router.POST("/lamp/all", getLampsHandler)
	router.POST("/lamp/add", addLampHandler)
	router.GET("/ping", func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "pong")
	})
	fmt.Println("Up and running")
	http.ListenAndServe(":10000", router)
}
