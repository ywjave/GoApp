package mysql

import (
	"GoApp/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) (err error) {
	sqlstr := `insert into post(
        	post_id,title,content,author_id,community_id)
			values (?,?,?,?,?)`
	_, err = db.Exec(sqlstr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}
func GetPostByID(pid int64) (post *models.Post, err error) {
	sqlstr := `select 
    post_id,title,content,author_id,community_id,create_time
	from post 
	where post_id = ? `
	post = new(models.Post)
	err = db.Get(post, sqlstr, pid)
	return
}

func GetPostList(page, size int64) (post []*models.Post, err error) {
	sqlstr := `select 
    post_id,title,content,author_id,community_id,create_time
	from post 
	order by create_time
	desc 
	limit ?,? `
	post = make([]*models.Post, 0, 2)
	err = db.Select(&post, sqlstr, (page-1)*size, size)
	return
}

func GetPostListByIDS(ids []string) (post []*models.Post, err error) {
	strsql := `select
	post_id,title,content,author_id,community_id,create_time
	from post 
	where post_id in (?)
	order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(strsql, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	err = db.Select(&post, query, args...)
	return
}
