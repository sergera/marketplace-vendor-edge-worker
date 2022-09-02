package evt

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sergera/marketplace-vendor-edge-worker/internal/domain"
	"github.com/sergera/marketplace-vendor-edge-worker/internal/service"
)

type OrderListener struct {
	service      *service.VendorAPIService
	eventHandler *EventHandler
}

func NewOrderListener() *OrderListener {
	return &OrderListener{service.NewVendorAPIService(), NewEventHandler()}
}

func (l *OrderListener) Listen() {
	l.eventHandler.Consume([]string{Topics[domain.Unconfirmed]}, func(msg kafka.Message) {
		newOrder := domain.OrderModel{}
		if err := json.Unmarshal(msg.Value, &newOrder); err != nil {
			log.Println("error unmarshalling kafka message: ", err.Error())
			return
		}
		l.NotifyVendor(newOrder)
	})
}

func (l *OrderListener) NotifyVendor(o domain.OrderModel) {
	l.service.SendOrder(o)
}
