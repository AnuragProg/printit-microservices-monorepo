import { z } from 'zod';

//type OrderEvent struct {
//	ShopId string `json:"shop_id"`
//	Status data.OrderStatus `json:"status"`
//}
const OrderStatusEnum = z.enum([
	"placed",
	"cancelled",
	"accepted",
	"rejected",
	"processing",
	"completed"
]);
const OrderEventSchema = z.object({
	shop_id: z.string(),
	status: OrderStatusEnum,
});
type OrderEvent = z.infer<typeof OrderEventSchema>;


export {
	OrderEvent,
	OrderStatusEnum,
	OrderEventSchema,
};
