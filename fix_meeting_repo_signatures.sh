#!/usr/bin/env bash
set -e
ROOT="$(pwd)"
echo "Running meeting-repo fix script from: $ROOT"

# 1) Update Meeting and MeetingMinutes model: RecordedBy -> uint
cat > backend/internal/models/stubs_generated.go <<'EOF'
package models

import "time"

// Meeting & related minimal models (used to satisfy compile-time refs)
type Meeting struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	MeetingType string    `json:"meeting_type,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at,omitempty"`
	Location    string    `json:"location,omitempty"`
	Agenda      string    `json:"agenda,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MeetingMinutes struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MeetingID  uint      `json:"meeting_id"`
	Attendees  string    `json:"attendees,omitempty"`
	Decisions  string    `json:"decisions,omitempty"`
	RecordedBy uint      `json:"recorded_by,omitempty"` // changed to uint to match service usage
	Content    string    `json:"content,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
EOF
echo "Updated backend/internal/models/stubs_generated.go (RecordedBy -> uint)"

# 2) Write a MeetingRepository with the signatures meeting_service.go expects
cat > backend/internal/repository/interfaces_generated.go <<'EOF'
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
EOF

echo "Wrote backend/internal/repository/interfaces_generated.go with expected method signatures"

# 3) Run tidy and attempt local build
cd backend
echo "Running go mod tidy..."
go mod tidy

echo "Attempting local build (showing up to 120 lines of output)..."
CGO_ENABLED=0 GOOS=linux go build -o /tmp/gram-server ./cmd/server 2>&1 | sed -n '1,120p' || true

cd "$ROOT"
echo "Done. If build produced more errors, paste the first 120 lines here and I'll update the stubs accordingly."
echo "If build succeeded, run:"
echo "  docker compose build --no-cache backend"
echo "  docker compose up -d backend frontend"
echo "  docker compose logs -f backend --tail 200"
