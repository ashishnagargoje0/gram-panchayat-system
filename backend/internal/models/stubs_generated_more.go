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
