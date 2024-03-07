package models

type ParamsSignup struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamsVoteData struct {
	//userid从请求中获取当前用户
	PostID    string `json:"post_id" binding:"required"`              //帖子id，
	Direction int8   `json:"direction,string" binding:"oneof=-1 0 1"` //赞成or反对，1 or -1 or 0
}
type ParamsPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"` //可以为空
	Page        int64  `form:"page"`
	Size        int64  `form:"size"`
	Order       string `form:"order"`
}

type ParamsCommunityPostList struct {
	*ParamsPostList
}

const (
	OrderTime = "time"
	OrerScore = "score"
)
