/*
* @Time ： 2023-02-06 16:17
* @Auth ： 张齐林
* @File ：user_dto.go
* @IDE ：GoLand
 */
package response

import "gin_applocation/model"

// UserDto 对应响应成功的data字段
type UserDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

// ToUserDto DTO 就是数据传输对象(Data Transfer Object)的缩写;用于展示层与服务层之间的数据传输对象
func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}
