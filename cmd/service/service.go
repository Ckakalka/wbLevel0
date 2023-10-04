package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Ckakalka/wbLevel0/db/postgres"
	"github.com/Ckakalka/wbLevel0/models"
	"github.com/Ckakalka/wbLevel0/provider"
	"github.com/Ckakalka/wbLevel0/server"
)

func main() {
	orderCash := models.NewOrderCash()
	dbManager, err := postgres.NewManager()
	if err != nil {
		log.Println(err)
	}

	allOrders, err := dbManager.GetAllOrders()
	for _, order := range allOrders {
		orderCash.Store(order.Uid, order)
	}

	provider := provider.NewStan(orderCash, dbManager)
	if err := provider.Start(); err != nil {
		log.Println(err)
	}

	closeAll := func() {
		if err := provider.Stop(); err != nil {
			log.Println(err)
		}
		if err := dbManager.Close(); err != nil {
			log.Println(err)
		}
	}
	sigInt := make(chan os.Signal)
	signal.Notify(sigInt, os.Interrupt)
	server := server.NewHttp(":8080", orderCash)
	go func() {
		<-sigInt
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		closeAll()
		os.Exit(1)
	}()

	if err := server.Start(); err != http.ErrServerClosed {
		closeAll()
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
