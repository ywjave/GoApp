package redis

import (
	"GoApp/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func GetPostIDSInOrder(p *models.ParamsPostList) ([]string, error) {
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrerScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//确定索引起始
	return getIdsFromKey(key, p.Page, p.Size)
}

func getIdsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	//zrevrange,（order）从大到小查询
	return rdb.ZRevRange(key, start, end).Result()
}

// GetPostIDSVote根据ids查询赞成票数
func GetPostIDSVote(ids []string) (data []int64, err error) {
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		keys = append(keys, key)
	}
	pipeline := rdb.Pipeline()
	for _, key := range keys {
		pipeline.ZCount(key, "1", "1")
	}
	//使用pipeline一次发送多个数据
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDSVote按社区根据ids查询赞成票数
func GetCommunityPostIDSInOrder(p *models.ParamsPostList) ([]string, error) {
	//使用zinterstore把分区的帖子与与帖子分数的zset生成一个新的zset
	//针对新的zset按之前的逻辑取数据

	orderkey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrerScore {
		orderkey = getRedisKey(KeyPostScoreZSet)
	}
	//社区key
	ckey := getRedisKey(KeycommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))
	//利用缓存减少执行次数？

	//从redis获取id
	//1.根据用户请求中携带的order参数确定要查询的redis key
	//缓存key减少zinterstore
	key := orderkey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(key).Val() < 1 {
		//不存在
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, ckey, orderkey)
		pipeline.Expire(key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec()            //一次执行多个命令
		if err != nil {
			return nil, err
		}
	}
	//拿到key后取值
	return getIdsFromKey(key, p.Page, p.Size)

}
