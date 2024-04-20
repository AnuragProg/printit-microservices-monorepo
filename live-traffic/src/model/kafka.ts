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
	// updated_on_or_before_epoch_ms: z.number().refine(epochMS => {
	// 	console.log(`received epoch timestamp in ms = ${epochMS}`);
	// 	const date = new Date(epochMS);
	// 	return !isNaN(date.getTime())
	// }, { message: 'require epoch timestamp in ms' }).transform(confirmedRfc3339 => new Date(confirmedRfc3339))
	updated_on_or_before_epoch_ms: z.number()
});

type OrderEvent = z.infer<typeof OrderEventSchema>;

export {
	OrderEvent,
	OrderStatusEnum,
	OrderEventSchema,
};
