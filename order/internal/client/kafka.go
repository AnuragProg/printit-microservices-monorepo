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
	ShopId string `json:"shop_id"`
	Status data.OrderStatus `json:"status"`
	UpdatedOnOrBefore string `json:"updated_on_or_before"`
}

func (oee *OrderEventEmitter) EmitOrderEvent(orderEvent *OrderEvent) error {
	if orderEvent == nil {
		return errors.New("no order event found")
	}
	orderEventJson, err := json.Marshal(orderEvent)
	if err != nil{
		return err
	}
	log.Infof("Emitting on kafka: msg = %v\n", orderEventJson);
	return oee.emitterClient.EmitSync(consts.ORDER_EMITTER_GROUP_ID, orderEventJson)
}
