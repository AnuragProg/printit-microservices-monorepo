

local untracked_shops = {}

for i = 1, #KEYS do
	local traffic = redis.call('GET', KEYS[i])

	--- apparently redis returns false value in case of nil reply
	if not traffic then
		table.insert(untracked_shops, KEYS[i])
	end
end

return untracked_shops
