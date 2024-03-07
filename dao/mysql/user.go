package mysql

import (
	"GoApp/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
)

func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExit
	}
	return nil
}

// 向表中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//执行SQL语句入库
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	user.Password = Encry(user.Password)
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	//fmt.Println("flag2", err)
	return
}

func LoginUser(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id ,username,password from user where username=?`

	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExit
	}
	if err != nil {
		return err
	}
	if user.Password != Encry(oPassword) {
		fmt.Println(user.Password, Encry(oPassword))
		return ErrorInvalidPassword
	}
	return nil
}

const secret = "jiami"

func Encry(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

func GetUserByID(user_id int64) (user *models.User, err error) {
	sqlstr := `select username from user where user_id=?`
	user = new(models.User)
	err = db.Get(user, sqlstr, user_id)
	return
}
