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
	order_id: z.string(),
	shop_id: z.string(),
	customer_id: z.string(),
	status: OrderStatusEnum,
	updated_on_or_before_epoch_ms: z.number()
});

type OrderEvent = z.infer<typeof OrderEventSchema>;

export {
	OrderEvent,
	OrderStatusEnum,
	OrderEventSchema,
};
