--检查参数数量
if #KEYS ~= 2 or #ARGV ~= 1 then
	return -1
end

local item_key = KEYS[1] --投票计数
local user_key = KEYS[2] --操作记录
local user_id = ARGV[1] --用户id

if redis.call("SISMEMBER", user_key, user_id) == 1 then
	return 0
end

redis.call("SADD", user_key, user_id)
redis.call("INCR", item_key)
return 1