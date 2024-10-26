package user

import (
	"log/slog"
	"time"

	"github.com/RGITHackathonFall2024/auth/internal/server"
	"gorm.io/gorm"
)

type User struct {
	TelegramID int64 `gorm:"primaryKey"`

	FirstName  string
	LastName   string
	HomeTown   string
	University string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func ByID(s *server.Server, id int64) (*User, error) {
	var user User
	log := s.Log().WithGroup("user-by-id")

	if err := s.DB().Where("telegram_id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("No such user", slog.Int64("telegram_id", id))
			return nil, &ErrNoSuchUser{}
		}

		return nil, err
	}

	return &user, nil
}

func Create(s *server.Server, telegramID int64, firstName, lastName string) error {
	return s.DB().Create(&User{
		TelegramID: telegramID,
		FirstName:  firstName,
		LastName:   lastName,
	}).Error
}

func Edit(s *server.Server, user *User) error {
	return s.DB().Model(&User{}).Where("telegram_id = ?", user.TelegramID).Updates(user).Error
}
