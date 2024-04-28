package client

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"

	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
	"github.com/AnuragProg/printit-microservices-monorepo/internal/data"
)


type OrderEventEmitter struct {
	emitterClient *goka.Emitter
}

func NewOrderEventEmitter(
	brokers []string,
) (*OrderEventEmitter, error) {
	// uses snappy compression by default
	client, err := goka.NewEmitter(brokers, goka.Stream(consts.ORDER_EVENT_TOPIC), new(codec.Bytes))
	if err != nil {
		return nil, err
	}
	return &OrderEventEmitter{
		emitterClient: client,
	}, nil
}

type OrderEvent struct {
	OrderId string `json:"order_id"`
	ShopId string `json:"shop_id"`
	CustomerId string `json:"customer_id"`
	Status data.OrderStatus `json:"status"`
	UpdatedOnOrBeforeEpochMS int64 `json:"updated_on_or_before_epoch_ms"`
}

func (oee *OrderEventEmitter) EmitOrderEvent(orderEvent *OrderEvent) error {
	if orderEvent == nil {
		return errors.New("no order event found")
	}
	log.Infof("shopId = %v, status = %v, updatedonorbefore = %v" , orderEvent.ShopId, orderEvent.Status, orderEvent.UpdatedOnOrBeforeEpochMS)
	orderEventJson, err := json.Marshal(orderEvent)
	if err != nil{
		return err
	}
	log.Infof("Emitting on kafka: msg = %v\n", orderEventJson);
	return oee.emitterClient.EmitSync(consts.ORDER_EMITTER_GROUP_ID, orderEventJson)
}
