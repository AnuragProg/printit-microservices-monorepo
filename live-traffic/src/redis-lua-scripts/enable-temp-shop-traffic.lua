
--------- inputs start ---------

local temp_traffic_key = KEYS[1]
local temp_traffic_timestamp_key = KEYS[2]

--------- inputs end ---------

local temp_traffic = tonumber(redis.call('GET', temp_traffic_key))
local temp_traffic_timestamp = tonumber(redis.call('GET', temp_traffic_timestamp_key))

if temp_traffic ~= nil and temp_traffic_timestamp ~= nil then
	-- return timestamp if already enabled
	return temp_traffic_timestamp
end

local redis_time = redis.call('TIME')
local secs = tonumber(redis_time[1])
local milli = tonumber(redis_time[2])/1000 -- because redis_time[2] is in microseconds

temp_traffic_timestamp = secs*1000 + milli -- epoch milliseconds

redis.call('SET', temp_traffic_key, 0)
redis.call('SET', temp_traffic_timestamp_key, temp_traffic_timestamp)

return temp_traffic_timestamp
