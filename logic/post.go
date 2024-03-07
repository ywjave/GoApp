package logic

import (
	"GoApp/dao/mysql"
	"GoApp/dao/redis"
	"GoApp/models"
	"GoApp/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1.生成postid
	p.ID = snowflake.GenID()

	//2.保存到数据库

	err = mysql.CreatePost(p)
	if err != nil {
		return
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
	//返回
}

func GetPostByID(pid int64) (data *models.PostDetail, err error) {

	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed",
			zap.Int64("pid", pid), zap.Error(err))
		return
	}
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}
	data = &models.PostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}
	return
}

func GetPostList(page, size int64) (data []*models.PostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return
	}
	data = make([]*models.PostDetail, 0, len(posts))

	for _, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.PostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(p *models.ParamsPostList) (data []*models.PostDetail, err error) {
	//redis查询id列表
	ids, err := redis.GetPostIDSInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDSInOrder(p) return 0 data")
		return
	}
	//提前查询数据？
	votedata, err := redis.GetPostIDSVote(ids)
	if err != nil {
		return nil, err
	}
	posts, err := mysql.GetPostListByIDS(ids)
	if err != nil {
		return
	}

	for index, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.PostDetail{
			AuthorName:      user.Username,
			VoteNum:         votedata[index],
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamsPostList) (data []*models.PostDetail, err error) {

	//redis查询id列表
	ids, err := redis.GetCommunityPostIDSInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDSInOrder(p) return 0 data")
		return
	}
	//提前查询数据？
	votedata, err := redis.GetPostIDSVote(ids)
	if err != nil {
		return nil, err
	}
	posts, err := mysql.GetPostListByIDS(ids)
	if err != nil {
		return
	}

	for index, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.PostDetail{
			AuthorName:      user.Username,
			VoteNum:         votedata[index],
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return
}

// 将两个查询接口整合
func GetPostListNew(p *models.ParamsPostList) (data []*models.PostDetail, err error) {
	if p.CommunityID == -1 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
	}
	return
}
