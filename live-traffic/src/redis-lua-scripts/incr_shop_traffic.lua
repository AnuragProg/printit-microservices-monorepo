-- PURPOSE: TO INCREMENT TRAFFIC ONLY WHEN TRAFFIC HAS ALREADY BEEN SET FOR THE SHOP
-- requires key[shop_id]


local shopId = KEYS[1]
local curShopTraffic = redis.call('GET', shopId)

-- no traffic is present (the shop traffic needs to be set explicitly)
if curShopTraffic == nil then
	return nil
end

curShopTraffic = tonumber(curShopTraffic)

-- not an integer (shop traffic must be integer to be incremented)
if curShopTraffic == nil then
	return nil
end

return redis.call('INCR', shopId)
