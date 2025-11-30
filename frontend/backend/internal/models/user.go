// internal/models/user.go
package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	
	Email             string         `gorm:"uniqueIndex;not null" json:"email"`
	Password          string         `gorm:"not null" json:"-"`
	Role              string         `gorm:"default:'citizen'" json:"role"` // admin, citizen
	
	// Personal Information
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	PhoneNumber       string         `gorm:"uniqueIndex" json:"phone_number"`
	AadharNumber      string         `gorm:"uniqueIndex" json:"aadhar_number"`
	DateOfBirth       *time.Time     `json:"date_of_birth"`
	Gender            string         `json:"gender"`
	
	// Address
	Address           string         `json:"address"`
	Village           string         `json:"village"`
	Taluka            string         `json:"taluka"`
	District          string         `json:"district"`
	State             string         `json:"state"`
	Pincode           string         `json:"pincode"`
	
	// Account Status
	IsVerified        bool           `gorm:"default:false" json:"is_verified"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	ProfileImage      string         `json:"profile_image"`
	
	// OTP for verification
	OTP               string         `json:"-"`
	OTPExpiry         *time.Time     `json:"-"`
	
	// Relations
	Applications      []Application  `gorm:"foreignKey:UserID" json:"applications,omitempty"`
	Complaints        []Complaint    `gorm:"foreignKey:UserID" json:"complaints,omitempty"`
	Properties        []Property     `gorm:"foreignKey:OwnerID" json:"properties,omitempty"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}