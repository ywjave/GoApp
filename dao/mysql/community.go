package mysql

import (
	"GoApp/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	sqlstr := "select community_id,community_name from community"
	if err = db.Select(&data, sqlstr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// 根据id查询社区详情
func GetCommunityDetailByID(id int64) (detail *models.CommunityDetail, err error) {
	sqlstr := `select
    community_id,community_name ,introduction,create_time 
	from community
	where community_id=?`
	detail = new(models.CommunityDetail)
	if err := db.Get(detail, sqlstr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no communityDetail in db")
			err = ErrorInvalidID
		}
	}
	return detail, err
}
