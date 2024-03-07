package redis

//redis key
//注意使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix              = "GoApp:"
	KeyPostTimeZSet        = "post:time"   //zset,帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  //zset,帖子及投票分数
	KeyPostVotedZSetPrefix = "post:voted:" //zset记录用户及投票类型，参数是post id?
	KeycommunitySetPrefix  = "community:"  //保存每个分区下的帖子id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
