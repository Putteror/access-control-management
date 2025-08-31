package service

import (
	"fmt"
	"time"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
)

type TimeRecordService interface {
	ClockIn(peopleID uint) error
	ClockOut(peopleID uint) error
	GetTimeRecordsByPeopleID(peopleID uint) ([]model.TimeRecord, error)
}

type timerecordServiceImpl struct {
	repo repository.TimeRecordRepository
}

func NewTimeRecordService(repo repository.TimeRecordRepository) TimeRecordService {
	return &timerecordServiceImpl{repo: repo}
}

func (s *timerecordServiceImpl) ClockIn(peopleID uint) error {
	record := model.TimeRecord{
		PeopleID:    peopleID,
		ClockInTime: time.Now(),
	}
	return s.repo.Create(&record)
}

func (s *timerecordServiceImpl) ClockOut(peopleID uint) error {
	// สมมติว่าต้องการหาบันทึกเวลาล่าสุดเพื่ออัปเดต
	records, err := s.repo.FindAllByPeopleID(peopleID)
	if err != nil || len(records) == 0 {
		return fmt.Errorf("no active clock-in record found for this user")
	}

	latestRecord := records[len(records)-1]
	latestRecord.ClockOutTime = time.Now()

	return s.repo.Create(&latestRecord)
}

func (s *timerecordServiceImpl) GetTimeRecordsByPeopleID(peopleID uint) ([]model.TimeRecord, error) {
	return s.repo.FindAllByPeopleID(peopleID)
}
