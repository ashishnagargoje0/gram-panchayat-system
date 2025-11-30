package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"gram-panchayat/internal/models"
	"gram-panchayat/internal/repository"
)
type MeetingService struct {
	meetingRepo *repository.MeetingRepository
}

func NewMeetingService(meetingRepo *repository.MeetingRepository) *MeetingService {
	return &MeetingService{meetingRepo: meetingRepo}
}

func (s *MeetingService) GetMeetings(page, limit int) ([]models.Meeting, int64, error) {
	return s.meetingRepo.List(page, limit)
}

func (s *MeetingService) GetMeeting(meetingID uint) (*models.Meeting, error) {
	return s.meetingRepo.GetByID(meetingID)
}

func (s *MeetingService) CreateMeeting(req interface{}) (*models.Meeting, error) {
	r := req.(*struct {
		Title       string
		Description string
		MeetingType string
		ScheduledAt string
		Location    string
		Agenda      string
	})

	scheduledAt, err := time.Parse(time.RFC3339, r.ScheduledAt)
	if err != nil {
		return nil, errors.New("invalid scheduled time format")
	}

	meeting := &models.Meeting{
		Title:       r.Title,
		Description: r.Description,
		MeetingType: r.MeetingType,
		ScheduledAt: scheduledAt,
		Location:    r.Location,
		Agenda:      r.Agenda,
		Status:      "scheduled",
	}

	if err := s.meetingRepo.Create(meeting); err != nil {
		return nil, err
	}

	return meeting, nil
}

func (s *MeetingService) AddMinutes(meetingID, adminID uint, req interface{}) (*models.MeetingMinutes, error) {
	r := req.(*struct {
		Content   string
		Attendees string
		Decisions string
	})

	minutes := &models.MeetingMinutes{
		MeetingID:  meetingID,
		Content:    r.Content,
		Attendees:  r.Attendees,
		Decisions:  r.Decisions,
		RecordedBy: adminID,
	}

	if err := s.meetingRepo.CreateMinutes(minutes); err != nil {
		return nil, err
	}

	// Update meeting status
	s.meetingRepo.Update(meetingID, map[string]interface{}{"status": "completed"})

	return minutes, nil
}