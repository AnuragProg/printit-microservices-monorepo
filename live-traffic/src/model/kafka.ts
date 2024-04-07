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
	updated_on_or_before: z.string().refine(rfc3339 => {
		const date = new Date(rfc3339);
		return !isNaN(date.getTime())
	}, { message: 'Invalid rfc3339 date format' }).transform(confirmedRfc3339 => new Date(confirmedRfc3339))
});

type OrderEvent = z.infer<typeof OrderEventSchema>;

export {
	OrderEvent,
	OrderStatusEnum,
	OrderEventSchema,
};
