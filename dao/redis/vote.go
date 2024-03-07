package redis

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneweektimesec = 7 * 24 * 3600
	scorePervote   = 432 //每一票432分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested   = errors.New("重复投票")
)

func CreatePost(postid, communityID int64) error {
	pipeline := rdb.Pipeline()
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postid,
	})
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postid,
	})
	//把帖子id加入社区的set
	ckey := getRedisKey(KeycommunitySetPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ckey, postid)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userid, postid string, value float64) error {
	//判断投票时间限制（1周内）
	posttime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postid).Val()
	if float64(time.Now().Unix())-posttime > oneweektimesec {
		fmt.Println(float64(time.Now().Unix()), posttime)
		return ErrVoteTimeExpire
	}
	//更新帖子分数
	pipeline := rdb.Pipeline()

	ov := rdb.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postid), userid).Val()
	if ov == value {
		return ErrVoteRepested
	}
	diff := math.Abs(ov - value) //计算两次投票的差值
	var dir float64

	if ov < value {
		dir = 1
	} else {
		dir = -1
	}
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), scorePervote*diff*dir, postid).Result()

	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postid), userid).Result()
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postid), redis.Z{
			Score:  value,
			Member: userid,
		}).Result()
	}
	_, err := pipeline.Exec()
	return err
}
