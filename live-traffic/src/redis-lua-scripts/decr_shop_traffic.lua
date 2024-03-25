-- PURPOSE: TO DECREMENT TRAFFIC PREVENTING NEGATIVE VALUES
-- requires key[shop_id]


local shopId = KEYS[1]
local curShopTraffic = redis.call('GET', shopId)

-- no traffic is present
if curShopTraffic == nil then
	return nil
end

curShopTraffic = tonumber(curShopTraffic)

-- not an integer and value should be greater than 0
if curShopTraffic == nil or curShopTraffic <= 0 then
	return nil
end

return redis.call('DECR', shopId)

