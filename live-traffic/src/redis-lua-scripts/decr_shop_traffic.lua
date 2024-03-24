-- requires key[shop_id]


local shopId = KEYS[1]
local curShopTraffic = redis.call('GET', shopId)

-- no traffic is present
if curShopTraffic == nil then
	return nil
end

curShopTraffic = tonumber(curShopTraffic)

-- not an integer
if curShopTraffic == nil then
	return nil
end

return redis.call('INCR', shopId)
