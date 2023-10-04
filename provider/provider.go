package provider

import (
	"encoding/json"
	"log"

	"github.com/Ckakalka/wbLevel0/db"
	"github.com/Ckakalka/wbLevel0/models"
	"github.com/nats-io/stan.go"
)

type Stan struct {
	orderCash *models.OrderCash
	stanConn  stan.Conn
	subcriber stan.Subscription
	dbManager db.Manager
}

func NewStan(cash *models.OrderCash, dbman db.Manager) *Stan {
	return &Stan{
		orderCash: cash,
		dbManager: dbman,
	}
}

func (s Stan) Start() error {
	var err error
	s.stanConn, err = stan.Connect("test-cluster", "valeriyReader")
	if err != nil {
		return err
	}
	messageHandler := func(m *stan.Msg) {
		var order db.Order
		if err := json.Unmarshal(m.Data, &order); err != nil {
			log.Println(err)
			return
		}
		if _, ok := s.orderCash.Load(order.Uid); !ok {
			s.orderCash.Store(order.Uid, order)
			if err := s.dbManager.InsertOrder(order); err != nil {
				log.Println(err)
			}
		} else {
			log.Printf("error order with uid = %s already exists\n", order.Uid)
		}
	}
	s.subcriber, err = s.stanConn.Subscribe("orders", messageHandler, stan.DeliverAllAvailable())
	if err != nil {
		return err
	}
	return nil
}

func (s Stan) Stop() error {
	if err := s.subcriber.Unsubscribe(); err != nil {
		return err
	}
	if err := s.stanConn.Close(); err != nil {
		return err
	}
	return nil
}
