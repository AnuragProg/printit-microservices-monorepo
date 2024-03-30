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
	updated_on_or_before: z.string().refine(rfc3999 => {
		const date = new Date(rfc3999);
		return !isNaN(date.getTime())
	}, { message: 'Invalid rfc3999 date format' }).transform(confirmedRfc3999 => new Date(confirmedRfc3999))
});

type OrderEvent = z.infer<typeof OrderEventSchema>;

export {
	OrderEvent,
	OrderStatusEnum,
	OrderEventSchema,
};
