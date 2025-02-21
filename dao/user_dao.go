package dao

import (
	"cabbage-server/db"
	"cabbage-server/model"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreateAccount 创建用户数据库操作
func CreateAccount(username, password string, email string) (*model.User, error) {
	user := &model.User{
		UserId:   uuid.New(),
		Name:     username,
		Password: password,
		Email:    email,
	}
	err := db.DB.Model(&model.User{}).Create(user).Error
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

// 根据邮箱查找用户
func FindUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := db.DB.Model(&model.User{}).Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

// 根据用户名查找
func FindUserByName(name string) (*model.User, error) {
	user := &model.User{}
	err := db.DB.Model(&model.User{}).
		Where("name = ?", name).
		First(user).
		Error
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func getUserRelationshipList(userid int64, page, size int, ship int) ([]*model.User, error) {
	result := []*model.User{}
	expression := db.DB.Model(&model.UserFollows{}).
		Select("uid").
		Where("ship = ?", ship).
		Where("uid = ?", userid)
	err := db.DB.Model(&model.User{}).
		Where("user in (?)", expression).
		Limit(size).
		Offset((page - 1) * size).
		Find(&result).
		Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 用户关注列表
func UserFollows(userid int64, page, size int) ([]*model.User, error) {
	result, err := getUserRelationshipList(userid, page, size, 1)
	return result, err
}

// 用户
func GetUserBlackList(userid int64, page, size int) ([]*model.User, error) {
	result, err := getUserRelationshipList(userid, page, size, -1)
	return result, err
}

// 根据用户名模糊查找
func FuzzyMatchingByUserName(name string) ([]*model.User, error) {
	users := []*model.User{}
	err := db.DB.Model(&model.User{}).
		Where("name like ?", "%"+name+"%").
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 统计某月新增用户
func CountNewUserOfMonth(month int) ([]*model.Counts, error) {
	timedate := time.Now()
	var results []*model.Counts
	year, _, _ := timedate.Date()
	dateString := strings.Join([]string{fmt.Sprintf("%d", year), fmt.Sprintf("%d", month), "01"}, "-")
	nextMonth := strings.Join([]string{fmt.Sprintf("%d", year), fmt.Sprintf("%d", month+1), "01"}, "-")
	err := db.DB.Model(&model.User{}).
		Select("DATE(created_at) as `date`", "count(*) as counts").
		Where("created_at >= ? AND created_at < ?", dateString, nextMonth).
		Order("DATE(created_at)").
		Group("date(created_at)").
		Find(&results).Error
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

// 统计当天新增用户数
func CountNewUserOfToday() (int64, error) {
	var count int64
	err := db.DB.Model(&model.User{}).
		Where("DATE(created_at) = CURDATE()").
		Count(&count).
		Error
	if err != nil {
		return -1, err
	}
	return count, nil
}
