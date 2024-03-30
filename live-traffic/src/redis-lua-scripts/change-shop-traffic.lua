
------- input start -----------

local perm_traffic_key = KEYS[1]
local temp_traffic_key = KEYS[2]
local temp_traffic_timestamp_key = KEYS[3] -- should contain milliseconds

local incr_or_decr = tostring(ARGV[1]) -- incr | decr
local order_epoch_timestamp = tonumber(ARGV[2]) -- in milliseconds (timestamp of order status updation)

------- input end -----------

local traffic_change = 0
if incr_or_decr == 'incr' then
	traffic_change = 1
elseif incr_or_decr == 'decr' then
	traffic_change = -1
else
	return 'INVALID_TRAFFIC_CHANGE_OPERATION'
end

if order_epoch_timestamp == nil then
	return 'INVALID_ORDER_TIMESTAMP'
end

local perm_traffic = tonumber(redis.call('GET', perm_traffic_key))
local temp_traffic = tonumber(redis.call('GET', temp_traffic_key))


---@param traffic number
local function is_traffic_valid(traffic)
	return traffic >= 0
end

if perm_traffic ~= nil then
	-- perm traffic is set, apply the the traffic change and offset the temp traffic value

	temp_traffic = temp_traffic or 0

	local new_traffic = perm_traffic + traffic_change + temp_traffic
	if is_traffic_valid(new_traffic) then
		redis.call('SET', perm_traffic_key, new_traffic)
		redis.call('DEL', temp_traffic_key, temp_traffic_timestamp_key)
		return new_traffic
	end

	redis.call('DEL', perm_traffic_key, temp_traffic_key, temp_traffic_timestamp_key)
	return 'NEGATIVE_TRAFFIC'
end

if temp_traffic == nil then
	return 'TEMP_TRAFFIC_NOT_ENABLED'
end

local temp_traffic_timestamp = tonumber(redis.call('GET', temp_traffic_timestamp_key))
if temp_traffic_timestamp == nil then
	redis.call('DEL', temp_traffic_key, temp_traffic_timestamp_key)
	return 'TEMP_TRAFFIC_TIMESTAMP_NOT_SET'
end

if order_epoch_timestamp <= temp_traffic_timestamp then
	return 'ORDER_NOT_ADDED_TO_TEMP'
end

local new_temp_traffic = temp_traffic + traffic_change
redis.call('SET', temp_traffic_key, new_temp_traffic)
return 'TRAFFIC_ADDED_TO_TEMP'


