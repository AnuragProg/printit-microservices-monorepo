package order

import (
	"time"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/AnuragProg/printit-microservices-monorepo/internal/client"
	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
	"github.com/AnuragProg/printit-microservices-monorepo/internal/data"
	auth "github.com/AnuragProg/printit-microservices-monorepo/proto_gen/authentication"
)


func GetOrderActionHandler(
	orderCol *mongo.Collection,
	orderEventEmitter *client.OrderEventEmitter,
	customerAllowedStatuses []data.OrderStatus,
	shopkeeperAllowedStatuses []data.OrderStatus,
) fiber.Handler{
	return func(c *fiber.Ctx) error {
		// extracting request/local data
		userInfo := c.Locals(consts.USER_INFO_LOCAL).(*auth.User)
		requestedStatus := c.Query("status")
		orderId, err := primitive.ObjectIDFromHex(c.Params("orderId"))
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, "invalid order id")
		}
		shopId := c.Params("shopId")

		// get statuses to check for validity of the status passed by user
		var statuses []data.OrderStatus
		switch userInfo.GetUserType() {
		case auth.UserType_CUSTOMER:
			statuses = customerAllowedStatuses
		case auth.UserType_SHOPKEEPER:
			statuses = shopkeeperAllowedStatuses
		default:
			log.Info("unknown user type detected here: " + userInfo.GetUserType().String())
			return fiber.ErrInternalServerError
		}

		// verify validity of requested status
		requestedStatusEnum, err := data.GetStatusEnum(requestedStatus)
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		isStatusValid := false
		for _, status := range statuses {
			if status == *requestedStatusEnum {
				isStatusValid = true
				break
			}
		}
		if !isStatusValid { // probably the customer tried shopkeeper's status or vice-versa
			return fiber.NewError(fiber.StatusUnauthorized, "status update not authorized")
		}

		// get order in question
		orderFilter := bson.M{"_id": orderId, "shop_id": shopId}
		result := orderCol.FindOne(context.Background(), orderFilter)
		if err := result.Err(); err != nil{
			log.Error(err.Error())
			return fiber.ErrInternalServerError
		}
		var orderInfo data.Order
		if err := result.Decode(&orderInfo); err != nil {
			log.Error(err.Error())
			return fiber.NewError(fiber.StatusNotFound, "order not found")
		}

		// execute order according to the status
		currentOrderStatus, err := data.GetStatusEnum(orderInfo.Status)
		if err != nil {
			log.Error(err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "something went wrong")
		}
		switch *requestedStatusEnum {
		case data.ORDER_ACCEPTED:
			if err := acceptOrder(orderCol, orderId, currentOrderStatus); err != nil {
				return err
			}
		case data.ORDER_REJECTED:
			if err := rejectOrder(orderCol, orderId, currentOrderStatus); err != nil {
				return err
			}
		case data.ORDER_CANCELLED:
			if err := cancelOrder(orderCol, orderId, currentOrderStatus); err != nil {
				return err
			}
		case data.ORDER_COMPLETED:
			if err := completeOrder(orderCol, orderId, currentOrderStatus); err != nil {
				return err
			}
		default:
			log.Warn("invalid request detected " + requestedStatus)
			return fiber.NewError(fiber.StatusMethodNotAllowed, "requested status not allowed")
		}

		// log the order event to kafka
		orderEvent := client.OrderEvent{
			ShopId: orderInfo.ShopId,
			Status: *requestedStatusEnum,
			UpdatedOnOrBefore: orderInfo.UpdatedAt.Format(time.RFC3339),
		}
		if err := orderEventEmitter.EmitOrderEvent(&orderEvent); err != nil {
			log.Error(err.Error())
		}

		c.JSON(map[string]interface{}{
			"message": "order " + string(requestedStatus) + " successfully",
		})

		return nil
	}
}

/**********************BELOW ALL METHODS RETURN FIBER ERROR INSTEAD OF GENERIC ERROR *********************************/
/**FOR NOW ACCEPT, REJECT, ETC.. METHODS ARE CALLING UPATEORDERSTATUS SO THAT THEY THEMSELVES DECIDE THE NEXT STATE**/

func acceptOrder(
	orderCol		*mongo.Collection,
	orderId		primitive.ObjectID,
	currentOrderStatus *data.OrderStatus,
) error {

	// will only allow acceptance when order is placed
	if *currentOrderStatus != data.ORDER_PLACED {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "status update not allowed on this order")
	}

	// update the status
	fromStatus, toStatus := string(data.ORDER_PLACED), string(data.ORDER_PROCESSING)
	if err := updateOrderStatus(orderCol, orderId, fromStatus, toStatus); err != nil {
		return err
	}

	return nil
}

func rejectOrder(
	orderCol		*mongo.Collection,
	orderId		primitive.ObjectID,
	currentOrderStatus *data.OrderStatus,
) error {

	// will only allow acceptance when order is placed
	if *currentOrderStatus != data.ORDER_PLACED {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "status update not allowed on this order")
	}

	// update the status
	fromStatus, toStatus := string(data.ORDER_PLACED), string(data.ORDER_REJECTED)
	if err := updateOrderStatus(orderCol, orderId, fromStatus, toStatus); err != nil {
		return err
	}

	return nil
}


func completeOrder(
	orderCol		*mongo.Collection,
	orderId		primitive.ObjectID,
	currentOrderStatus *data.OrderStatus,
) error {

	// will only allow acceptance when order is placed
	if *currentOrderStatus != data.ORDER_PROCESSING {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "status update not allowed on this order")
	}

	// update the status
	fromStatus, toStatus := string(data.ORDER_PROCESSING), string(data.ORDER_COMPLETED)
	if err := updateOrderStatus(orderCol, orderId, fromStatus, toStatus); err != nil {
		return err
	}

	return nil
}

func cancelOrder(
	orderCol		*mongo.Collection,
	orderId		primitive.ObjectID,
	currentOrderStatus *data.OrderStatus,
) error {

	// will only allow acceptance when order is placed
	if *currentOrderStatus != data.ORDER_PLACED {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "status update not allowed on this order")
	}

	// update the status
	fromStatus, toStatus := string(data.ORDER_PLACED), string(data.ORDER_CANCELLED)
	if err := updateOrderStatus(orderCol, orderId, fromStatus, toStatus); err != nil {
		return err
	}
	return nil
}


/* This method directly changes status of order, hence proper guards needs to be called before it */
func updateOrderStatus(
	orderCol		*mongo.Collection,
	orderId		primitive.ObjectID,
	fromStatus  string, // from and to status prevents race conditions where undefined state transition may take place
	toStatus		string,
) error {
	orderFilter := bson.M{"_id": orderId, "status": fromStatus}
	orderUpdate := bson.M{
		"$set": bson.M{
			"status": toStatus,
			"updated_at": time.Now().UTC(),
		},
	}
	res, err := orderCol.UpdateOne(context.Background(), orderFilter, orderUpdate)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if res.ModifiedCount == 0 {
		return fiber.NewError(fiber.StatusInternalServerError, "status unable to update, something went wrong!")
	}
	return nil
}

