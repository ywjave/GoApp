package logic

import (
	"GoApp/dao/mysql"
	"GoApp/models"
)

func GetCommunityList() ([]*models.Community, error) {
	//数据库查询所有社区，并返回结果
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	//数据库查询所有社区，并返回结果
	return mysql.GetCommunityDetailByID(id)
}
