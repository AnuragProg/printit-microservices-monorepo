
--- input start

local perm_traffic_key = KEYS[1]
local temp_traffic_key = KEYS[2]
local temp_traffic_timestamp_key = KEYS[3] -- should contain milliseconds

--- input end

local perm_traffic = tonumber(redis.call('GET', perm_traffic_key))
if perm_traffic == nil then
	return 'NOT_FOUND'
end

--- update traffic on read (main strategy that we will be following)
local temp_traffic = tonumber(redis.call('GET', temp_traffic_key)) or 0
local new_traffic = perm_traffic + temp_traffic

if new_traffic < 0 then
	redis.call('DEL', perm_traffic_key, temp_traffic_key, temp_traffic_timestamp_key)
	return 'NEGATIVE_TRAFFIC'
end

redis.call('SET', perm_traffic_key, new_traffic)
redis.call('DEL', temp_traffic_key, temp_traffic_timestamp_key)

return new_traffic
