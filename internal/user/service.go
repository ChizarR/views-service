package user

import (
	"context"
	"time"

	"github.com/ChizarR/stats-service/pkg/logging"
)

type UserStatService interface {
	UpdateUserViews(ctx context.Context, uDTO UserDTO) error
	GetUserViewsForToday(ctx context.Context, tgId int) (UserDTO, error)
	GetAllUsersStats(ctx context.Context) ([]User, error)
}

type service struct {
	storage Storage
	logger  *logging.Logger
}

func NewUserStatService(storage Storage, logger *logging.Logger) UserStatService {
	return &service{storage: storage, logger: logger}
}

func (s *service) UpdateUserViews(ctx context.Context, uDTO UserDTO) error {
	today := getTodayDate()
	user, err := s.storage.GetOrCreate(ctx, uDTO.TgId)
	if err != nil {
		return err
	}

	var dateFound bool
	if len(user.Intaractions) != 0 {
		for idx, day := range user.Intaractions {
			if day.Date == today {
				dateFound = true
				for category, viewsNum := range uDTO.Views {
					_, ok := day.Views[category]
					if !ok {
						day.Views[category] = viewsNum
					}
					day.Views[category] += viewsNum
				}
				user.Intaractions[idx] = day
				break
			}
		}
		if !dateFound {
			todayIntaractions := Intaractions{
				Date:  today,
				Views: uDTO.Views,
			}
			user.Intaractions = append(user.Intaractions, todayIntaractions)
		}
	} else {
		todayIntaractions := Intaractions{
			Date:  today,
			Views: uDTO.Views,
		}
		user.Intaractions = append(user.Intaractions, todayIntaractions)
	}

	if err := s.storage.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *service) GetUserViewsForToday(ctx context.Context, tgId int) (UserDTO, error) {
	today := getTodayDate()
	user, err := s.storage.GetOrCreate(ctx, tgId)
	if err != nil {
		return UserDTO{}, err
	}

	uDTO := UserDTO{
		TgId:  user.TgId,
		Views: map[string]int{},
	}
	intaractions := user.Intaractions
	for _, day := range intaractions {
		if day.Date == today {
			uDTO.Views = day.Views
		}
	}

	return uDTO, nil
}

func (s *service) GetAllUsersStats(ctx context.Context) ([]User, error) {
	allUsers, err := s.storage.FindAll(ctx)
	if err != nil {
		return []User{}, err
	}
	return allUsers, nil
}

func getTodayDate() string {
	const YYYYMMDD = "2006-01-02"
	t := time.Now().UTC()
	d := t.Format(YYYYMMDD)
	return d
}
