package models

import (
	database "example/auth/internal/pkg/database"
	"time"

	"gorm.io/gorm"
)

type Org struct {
	Id   uint `gorm:"primaryKey;column:id"`
	Name string
}

type User struct {
	Username string `gorm:"primaryKey"`
	Password string `gorm:"not null"`
	Isadmin  bool   `gorm:"not null; column:admin"`
	OrgId    uint   `gorm:"column:orgid"`
	Org      Org    `gorm:"foreignKey:OrgId"`
}

type JwtRefreshToken struct {
	Username     string `gorm:"primaryKey"`
	RefreshToken string `gorm:"column:refreshToken"`
	Expire       int64  `gorm:"column:expire"`
}

type JwtBlackListedToken struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	AccessToken string `gorm:"column:accessToken"`
}

// delete the blacklisted token after one hour (access token validity) since it will be invalid after that anyway a user won't be able to use it
func (at *JwtBlackListedToken) AfterCreate(db *gorm.DB) (err error) {
	time.AfterFunc(time.Hour, func() {
		db.Delete(at)
	})
	return
}

func init() {
	// create database table for org and user and refreshToken
	database.Db.AutoMigrate(&Org{})
	database.Db.AutoMigrate(&User{})
	database.Db.AutoMigrate(&JwtRefreshToken{})
	database.Db.AutoMigrate(&JwtBlackListedToken{})
}
