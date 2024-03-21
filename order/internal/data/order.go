package data

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	ORDER_PLACED		OrderStatus = "placed"
	ORDER_CANCELLED	OrderStatus = "cancelled"
	ORDER_ACCEPTED		OrderStatus = "accepted"
	ORDER_REJECTED		OrderStatus = "rejected"
	ORDER_PROCESSING	OrderStatus = "processing"
	ORDER_COMPLETED	OrderStatus = "completed"
)

type Order struct {
	Id					primitive.ObjectID `bson:"_id" json:"_id"`
	FileId			string `bson:"file_id" json:"file_id"`
	ShopId			string `bson:"shop_id" json:"shop_id"`
	PriceId			string `bson:"price_id" json:"price_Id"`
	CustomerId		string `bson:"customer_id" json:"customer_id"`
	Status 			string `bson:"status" json:"status"`
	CreatedAt		time.Time `bson:"created_at" json:"-"`
	UpdatedAt		time.Time `bson:"updated_at" json:"-"`
}

func GetStatusEnum(status string) (*OrderStatus, error) {
	validStatuses := []OrderStatus{ORDER_PLACED, ORDER_CANCELLED, ORDER_ACCEPTED, ORDER_REJECTED, ORDER_PROCESSING, ORDER_COMPLETED}

	for _, validStatus := range validStatuses {
		if string(validStatus) == status {
			return &validStatus, nil
		}
	}
	return nil, errors.New("invalid status")
}


func CreateOrder(fileId, shopId, priceId, customerId, status string) (*Order, error) {
	statusEnum, err := GetStatusEnum(status)
	if err != nil{
		return nil, err
	}
	return &Order{
		Id: primitive.NewObjectID(),
		FileId: fileId,
		ShopId: shopId,
		PriceId: priceId,
		CustomerId: customerId,
		Status: string(*statusEnum),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
