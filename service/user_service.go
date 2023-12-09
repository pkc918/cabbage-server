package service

import (
	"cabbage-server/dao"
	"cabbage-server/dto"
	"cabbage-server/internal"
	"cabbage-server/model"
	"errors"

	"gorm.io/gorm"
)

// CreateAccount 创建新用户服务
func CreateAccount(user *dto.SignupDTO) error {
	err := dao.CreateAccount(&model.User{
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return internal.UserNotFoundError
		} else {
			return internal.InernalError
		}
	}
	return nil
}

// GetUserProfile 获取用户信息服务
func GetUserProfile(email string) (*model.User, error) {
	user, err := dao.FindUserByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internal.RecordNotFoundError
		} else {
			return nil, internal.InernalError
		}
	}
	return user, nil
}

func FindUserByName(name string) (*model.User, error) {
	user, err := dao.FindUserByName(name)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internal.InernalError
		}
		return nil, internal.UserNotFoundError
	}
	return user, nil
}

func ProfileShare(user string) (*model.UserProfile, error) {
	_user, err := dao.FindUserByName(user)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internal.InernalError
		}
		return nil, internal.RecordNotFoundError
	}
	profile, err := dao.FindProfileByUID(int64(_user.ID))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internal.InernalError
		}
		return nil, internal.RecordNotFoundError
	}
	return profile, nil
}
