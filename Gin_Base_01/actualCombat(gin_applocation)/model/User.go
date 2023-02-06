/*
* @Time ： 2023-02-06 15:38
* @Auth ： 张齐林
* @File ：User.go
* @IDE ：GoLand
 */
package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}
