package logic

import (
	"GoApp/dao/mysql"
	"GoApp/models"
	"GoApp/pkg/jwt"
	"GoApp/pkg/snowflake"
)

func SignUp(p *models.ParamsSignup) (err error) {
	//判断用户是否存在

	if err = mysql.CheckUserExist(p.Username); err != nil {

		//fmt.Println("logic", err)
		return err
	}
	//生成uid
	userID := snowflake.GenID()
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.密码加密
	//保存到数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamsLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//user是指针
	if err := mysql.LoginUser(user); err != nil {
		return nil, err
	}
	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
