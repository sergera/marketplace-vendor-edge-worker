package api

import (
	"encoding/json"
	"net/http"

	"github.com/sergera/marketplace-vendor-edge-worker/internal/domain"
	"github.com/sergera/marketplace-vendor-edge-worker/internal/evt"
)

type OrderAPI struct {
	eventHandler *evt.EventHandler
}

func NewOrderAPI() *OrderAPI {
	return &OrderAPI{evt.NewEventHandler()}
}

func (o *OrderAPI) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var m domain.OrderModel

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderInBytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status, err := m.StatusType()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(orderInBytes)

	o.eventHandler.Produce(evt.Topics[status], "", orderInBytes)
}
