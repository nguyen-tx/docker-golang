package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"     // Quản trị viên hệ thống
	RolePilot     UserRole = "pilot"     // Phi công / Người điều khiển UAV
	RoleOperator  UserRole = "operator"  // Đơn vị vận hành
	RoleATC       UserRole = "atc"       // Kiểm soát không lưu
	RoleInspector UserRole = "inspector" // Thanh tra hàng không
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending" // Chờ xét duyệt
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	FullName     string     `db:"full_name" json:"full_name"`
	Phone        string     `db:"phone" json:"phone"`
	Role         UserRole   `db:"role" json:"role"`
	Status       UserStatus `db:"status" json:"status"`
	LicenseNo    string     `db:"license_no" json:"license_no,omitempty"` // Số giấy phép bay
	Organization string     `db:"organization" json:"organization,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

type UserProfile struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	FullName     string     `json:"full_name"`
	Phone        string     `json:"phone"`
	Role         UserRole   `json:"role"`
	Status       UserStatus `json:"status"`
	LicenseNo    string     `json:"license_no,omitempty"`
	Organization string     `json:"organization,omitempty"`
}
