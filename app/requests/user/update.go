package user

import (
	userModel "gin_bbs/app/models/user"
	"gin_bbs/app/requests"
	"gin_bbs/pkg/ginutils/flash"
	"gin_bbs/pkg/ginutils/validate"

	"github.com/gin-gonic/gin"
)

type UserUpdateForm struct {
	validate.Validate
	ID           int
	Name         string
	Email        string
	Introduction string
}

// IsStrict 有错误即退出
func (*UserUpdateForm) IsStrict() bool {
	return true
}

// RegisterValidators 注册验证器
func (u *UserUpdateForm) RegisterValidators() validate.ValidatorMap {
	return validate.ValidatorMap{
		"name": {
			validate.RequiredValidator(u.Name),
			validate.BetweenValidator(u.Name, 3, 25),
			validate.RegexpValidator(u.Name, `^[A-Za-z0-9\-\_]+$`),
			requests.NameUniqueValidator(u.Name, u.ID),
		},
		"email": {
			validate.RequiredValidator(u.Email),
			validate.MaxLengthValidator(u.Email, 255),
			validate.EmailValidator(u.Email),
		},
		"introduction": {
			validate.MaxLengthValidator(u.Introduction, 80),
		},
	}
}

func (*UserUpdateForm) RegisterMessages() validate.MessagesMap {
	return validate.MessagesMap{
		"name": {
			"用户名不能为空。",
			"用户名必须介于 3 - 25 个字符之间",
			"用户名只支持英文、数字、横杠和下划线。",
			"用户名已被占用，请重新填写",
		},
	}
}

// ValidateAndUpdate 验证参数并且更新用户
func (u *UserUpdateForm) ValidateAndUpdate(c *gin.Context, user *userModel.User) bool {
	ok, errArr, errMap := validate.Run(u)
	if !ok {
		validate.SaveValidateMessage(c, errArr, errMap)
		return false
	}

	// 更新用户
	user.Name = u.Name
	user.Email = u.Email
	user.Introduction = u.Introduction

	if err := user.Update(); err != nil {
		flash.NewDangerFlash(c, "用户更新失败: "+err.Error())
		return false
	}

	return true
}