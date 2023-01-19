package intaraction

import (
	"context"
	"time"

	"github.com/ChizarR/stats-service/pkg/logging"
)

type IntaractionService interface {
	AddNewIntaractions(ctx context.Context, intrDTO IntaractionDTO) error
	GetTodayIntaractions(ctx context.Context) (IntaractionDTO, error)
	GetAllIntaractions(ctx context.Context) (IntaractionDTO, error)
}

type service struct {
	storage Storage
	logger  *logging.Logger
}

func NewIntaractionService(storage Storage, logger *logging.Logger) IntaractionService {
	return &service{storage: storage, logger: logger}
}

func (s *service) AddNewIntaractions(ctx context.Context, intrDTO IntaractionDTO) error {
	date := getTodayDate()
	todayStats, err := s.storage.GetOrCreate(ctx, date)
	if err != nil {
		return err
	}

	for key, value := range intrDTO.Views {
		_, ok := todayStats.Views[key]
		if !ok {
			todayStats.Views[key] = value
			continue
		}
		todayStats.Views[key] += value
	}

	if err := s.storage.Update(ctx, todayStats); err != nil {
		return err
	}
	return nil
}

func (s *service) GetTodayIntaractions(ctx context.Context) (IntaractionDTO, error) {
	date := getTodayDate()
	todayStats, err := s.storage.GetOrCreate(ctx, date)
	if err != nil {
		return IntaractionDTO{}, err
	}
	intrDTO := IntaractionDTO{
		Views: todayStats.Views,
	}
	return intrDTO, nil
}

func (s *service) GetAllIntaractions(ctx context.Context) (IntaractionDTO, error) {
	allIntrs, err := s.storage.FindAll(ctx)
	if err != nil {
		return IntaractionDTO{}, err
	}

	result := IntaractionDTO{Views: map[string]int{}}
	for _, intr := range allIntrs {
		for key, value := range intr.Views {
			_, ok := result.Views[key]
			if !ok {
				result.Views[key] = value
				continue
			}
			result.Views[key] += value
		}
	}
	return result, nil
}

func getTodayDate() string {
	const YYYYMMDD = "2006-01-02"
	t := time.Now().UTC()
	d := t.Format(YYYYMMDD)
	return d
}
