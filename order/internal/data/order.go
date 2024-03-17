package data

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


const (
	ORDER_PLACED = "placed"
	ORDER_CANCELLED= "cancelled"
	ORDER_ACCEPTED = "accepted"
	ORDER_REJECTED = "rejected"
	ORDER_PROCESSING = "processing"
	ORDER_COMPLETED = "completed"
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

func IsStatusValid(status string) error {
	validStatuses := []string{ORDER_PLACED, ORDER_CANCELLED, ORDER_ACCEPTED, ORDER_REJECTED, ORDER_PROCESSING, ORDER_COMPLETED}
	isValid := false

	for _, validStatus := range validStatuses {
		if validStatus == status {
			isValid = true
			break
		}
	}

	if isValid {
		return nil
	}

	return errors.New("invalid status")
}


func CreateOrder(fileId, shopId, priceId, customerId, status string) (*Order, error) {
	if err := IsStatusValid(status); err != nil{
		return nil, err
	}
	return &Order{
		Id: primitive.NewObjectID(),
		FileId: fileId,
		ShopId: shopId,
		PriceId: priceId,
		CustomerId: customerId,
		Status: status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
