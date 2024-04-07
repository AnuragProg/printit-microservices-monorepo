import { z } from 'zod';

const GetShopTrafficSchema = z.object({
	action: z.enum(['get']),
	shopIds: z.array(z.string()),
});

type GetShopTraffic = z.infer<typeof GetShopTrafficSchema>;

const SubscriptionShopTrafficSchema = z.object({
	action: z.enum(['subscribe', 'unsubscribe']),
	shopIds: z.array(z.string()),
});

type SubscriptionShopTraffic = z.infer<typeof SubscriptionShopTrafficSchema>;

export {
	GetShopTraffic,
	GetShopTrafficSchema,
	SubscriptionShopTraffic,
	SubscriptionShopTrafficSchema,
};
