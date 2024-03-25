import { z } from 'zod';


const SubscriptionShopTrafficSchema = z.object({
	action: z.enum(['subscribe', 'unsubscribe']),
	shopIds: z.array(z.string()),
});

type SubscriptionShopTraffic = z.infer<typeof SubscriptionShopTrafficSchema>;

export {
	SubscriptionShopTraffic,
	SubscriptionShopTrafficSchema,
};
