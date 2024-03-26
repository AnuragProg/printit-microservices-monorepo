local temp_key = KEYS[1]
local perm_key = KEYS[2]

redis.log(redis.LOG_DEBUG, 'temp_key = '..temp_key)
redis.log(redis.LOG_DEBUG, 'perm_key = '..perm_key)

local temp_traffic = tonumber(redis.call('GET', temp_key))
local perm_traffic = tonumber(redis.call('GET', perm_key))

if perm_traffic == nil then
	-- shop is not being tracked
	if temp_traffic == nil then
		-- temp is not created for storing temp traffic yet
		redis.log(redis.LOG_WARNING, 'not tracking shop with temp_key = '..temp_key..' and perm_key = '..perm_key)
		return redis.error_reply('UNTRACKED_TEMP_TRAFFIC')
	end
	redis.call('INCR', temp_key)
	return redis.error_reply('UNTRACKED_PERM_TRAFFIC')
end

temp_traffic = temp_traffic or 0
local new_traffic = perm_traffic + 1 + temp_traffic
if new_traffic < 0 then
	redis.log(redis.LOG_WARNING, 'negative traffic: encountered with value '..(new_traffic))
	redis.log(redis.LOG_WARNING, 'negative traffic: deleting temp and perm shop traffic')
	redis.call('DEL', temp_key, perm_key)
	return redis.error_reply('INVALID_TRAFFIC')
end

redis.call('SET', perm_key, new_traffic)
redis.call('DEL', temp_key)

return new_traffic
