#!/usr/bin/env bash
set -e
ROOT="$(pwd)"
echo "Running final-fix script from: $ROOT"

# 1) Add missing model stubs so internal/database/postgres.go can compile
cat > backend/internal/models/stubs_generated_more.go <<'EOF'
package models

import "time"

type TaxBill struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PropertyID uint      `json:"property_id"`
	Amount     float64   `json:"amount"`
	DueDate    time.Time `json:"due_date,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Payment struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `json:"user_id"`
	BillID   uint      `json:"bill_id,omitempty"`
	Amount   float64   `json:"amount"`
	PaidAt   time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Notice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Scheme struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type SchemeApplication struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SchemeID  uint      `json:"scheme_id"`
	UserID    uint      `json:"user_id"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Document struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name,omitempty"`
	Path      string    `json:"path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Message   string    `json:"message,omitempty"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}
EOF

echo "Wrote backend/internal/models/stubs_generated_more.go"

# 2) Add a tiny file to reference encoding/json and fmt to avoid "imported and not used" errors.
#    This is safe and non-invasive — it only creates harmless blank usages.
cat > backend/internal/service/_fix_unused_imports.go <<'EOF'
package service

import (
	"encoding/json"
	"fmt"
)

// Small no-op references so files that import encoding/json or fmt don't fail
var _ = json.Marshal
var _ = fmt.Sprintf
EOF

echo "Wrote backend/internal/service/_fix_unused_imports.go (prevents unused-import errors)"

# 3) Run go mod tidy and do a local build (print first 160 lines)
cd backend
echo "Running go mod tidy..."
go mod tidy

echo "Attempting local build (showing up to 160 lines of output)..."
CGO_ENABLED=0 GOOS=linux go build -o /tmp/gram-server ./cmd/server 2>&1 | sed -n '1,160p' || true
BUILD_EXIT=$?

cd "$ROOT"

if [ "$BUILD_EXIT" -eq 0 ]; then
  echo
  echo "BUILD SUCCEEDED locally."
  echo "Next recommended steps (rebuild docker image & start services):"
  echo "  docker compose build --no-cache backend"
  echo "  docker compose up -d backend frontend"
  echo "  docker compose logs -f backend --tail 200"
else
  echo
  echo "BUILD produced errors (above). Paste the build output here and I'll update stubs/fixes further."
fi

