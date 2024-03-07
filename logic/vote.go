package logic

import (
	"GoApp/dao/redis"
	"GoApp/models"
	"strconv"

	"go.uber.org/zap"
)

//投票功能 www.ruanyifeng.com/blog/algorithm/

//使用简化算法
//投一票加432分？

func VoteForPost(userid int64, p *models.ParamsVoteData) error {
	zap.L().Debug("voteForPost",
		zap.Int64("userid", userid),
		zap.String("postid", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userid)), p.PostID, float64(p.Direction))
}
