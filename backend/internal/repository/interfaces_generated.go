package repository

import "gram-panchayat/internal/models"

// Concrete minimal MeetingRepository with signatures used by services.
// These are no-op implementations — replace with real DB logic later.

type MeetingRepository struct{}

// List accepts pagination and returns slice, totalCount and error
func (r *MeetingRepository) List(limit int, offset int) ([]models.Meeting, int64, error) {
	return []models.Meeting{}, 0, nil
}

// GetByID returns a pointer to Meeting and error
func (r *MeetingRepository) GetByID(id uint) (*models.Meeting, error) {
	m := &models.Meeting{}
	return m, nil
}

// Create inserts a meeting (no-op)
func (r *MeetingRepository) Create(m *models.Meeting) error {
	return nil
}

// Update accepts id and fields map (as used by service)
func (r *MeetingRepository) Update(id uint, fields map[string]interface{}) error {
	return nil
}

// CreateMinutes is a no-op stub
func (r *MeetingRepository) CreateMinutes(mm *models.MeetingMinutes) error {
	return nil
}

// Keep other repository names as simple interfaces for now
type UserRepository interface{}
type PropertyRepository interface{}
type ComplaintRepository interface{}
type ApplicationRepository interface{}
type NoticeRepository interface{}
type SchemeRepository interface{}
type PaymentRepository interface{}
type DocumentRepository interface{}
type NotificationRepository interface{}
