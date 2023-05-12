local key = "rate:d:leaky:bucket:" .. KEYS[1]
local peak = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local is_exists = redis.call("EXISTS", key)
if is_exists == 0 then
    redis.call("HMSET", key, "last_time", now, "peak", peak, "rate", rate, "cur", 0)
end

local bucket = redis.pcall("HMGET", key, "last_time", "peak", "rate", "cur")
local bucket_last_time = tonumber(bucket[1])
local bucket_peak = tonumber(bucket[2])
local bucket_rate = tonumber(bucket[3])
local bucket_cur = tonumber(bucket[4])

if rate ~= bucket_rate or peak ~= bucket_peak then
    bucket_rate = rate
    bucket_peak = peak
    redis.pcall("HMSET", key, "peak", peak, "rate", rate)
end

local cur = bucket_cur - (now-bucket_last_time)*bucket_rate
if cur < 0 then
    cur = 0
end

if cur >= bucket_peak then
    redis.pcall("HMSET", key, "last_time", now)
    return 0
end

redis.pcall("HMSET", key, "last_time", now, "cur", cur+1)

return 1