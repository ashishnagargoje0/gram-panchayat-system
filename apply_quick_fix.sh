#!/usr/bin/env bash
set -e

ROOT="$(pwd)"
echo "Running quick-fix script from: $ROOT"

# 1) Fix backend/.env (replace the corrupted file)
cat > backend/.env <<'ENV'
# Backend environment for local/dev
POSTGRES_HOST=gp_db
POSTGRES_PORT=5432
POSTGRES_USER=gp_user
POSTGRES_PASSWORD=gp_password
POSTGRES_DB=gram_panchayat
APP_PORT=8080
JWT_SECRET=replace_me_with_real_secret
ENV
echo "Wrote backend/.env"

# 2) Update models: add RecordedBy to MeetingMinutes (and keep existing stubs)
cat > backend/internal/models/stubs_generated.go <<'EOF'
package models

import "time"

// Minimal meeting-related models to satisfy compile-time references.
// Replace these with real models later.

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
	RecordedBy string    `json:"recorded_by,omitempty"` // added to match code references
	Content    string    `json:"content,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
EOF
echo "Updated Meeting/MeetingMinutes model stub"

# 3) Replace repository stubs: implement a minimal concrete MeetingRepository so *repository.MeetingRepository has methods
cat > backend/internal/repository/interfaces_generated.go <<'EOF'
package repository

import "gram-panchayat/internal/models"

// Minimal repository definitions to unblock compilation.
// The MeetingRepository below is a concrete struct with no-op methods
// to satisfy calls like s.meetingRepo.List(...) where s.meetingRepo is a pointer.

type MeetingRepository struct{}

// List returns an empty slice and nil error
func (r *MeetingRepository) List() ([]models.Meeting, error) {
	return []models.Meeting{}, nil
}

// GetByID returns zero Meeting and error nil
func (r *MeetingRepository) GetByID(id uint) (models.Meeting, error) {
	return models.Meeting{}, nil
}

// Create is a no-op stub
func (r *MeetingRepository) Create(m *models.Meeting) error {
	return nil
}

// Update is a no-op stub
func (r *MeetingRepository) Update(m *models.Meeting) error {
	return nil
}

// CreateMinutes is a no-op stub
func (r *MeetingRepository) CreateMinutes(mm *models.MeetingMinutes) error {
	return nil
}

// Keep other repository names as simple interfaces (replace later with real methods)
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
echo "Wrote repository stubs (concrete MeetingRepository implemented)"

# 4) Run module tidy and attempt build
cd backend
echo "Running go mod tidy..."
go mod tidy

echo "Attempting local build (showing up to 80 lines of output)..."
CGO_ENABLED=0 GOOS=linux go build -o /tmp/gram-server ./cmd/server 2>&1 | sed -n '1,80p' || true

echo "Done. If build succeeded, you can rebuild Docker images next."
cd "$ROOT"

echo "Suggested next commands if build is successful:"
echo "  docker compose build --no-cache backend"
echo "  docker compose up -d backend frontend"
echo "  docker compose logs -f backend --tail 200"
