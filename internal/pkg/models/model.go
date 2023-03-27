package models

import (
	database "example/auth/internal/pkg/database"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username string `gorm:"primaryKey"`
	Password string `gorm:"not null"`
	Admin    bool   `gorm:"not null; column:admin"`
	OrgId    uint   `gorm:"column:orgid"`
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
	// Drop table if exist i.e Delete JwtBlackListedToken on server restart
	// database.Db.Migrator().DropTable(&JwtRefreshToken{})
	database.Db.Migrator().DropTable(&JwtBlackListedToken{})

	// create database table
	database.Db.AutoMigrate(&User{})
	database.Db.AutoMigrate(&JwtRefreshToken{})
	database.Db.AutoMigrate(&JwtBlackListedToken{})
}
