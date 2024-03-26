local perm_key = KEYS[1]
local new_traffic = tonumber(ARGV[1])

redis.log(redis.LOG_DEBUG, 'perm_key = '..perm_key)
redis.log(redis.LOG_DEBUG, 'new_traffic = '..new_traffic)

if new_traffic == nil then
	return redis.error_reply('INVALID_ARGUMENT')
end

redis.call('SET', perm_key, new_traffic)
