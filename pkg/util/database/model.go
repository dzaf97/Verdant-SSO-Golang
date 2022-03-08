package database

import (
	"time"
)

//DATABASE TABLE
type (
	GormDefault struct {
		CreatedAt time.Time  `json:"createdAt" `
		UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
		DeletedAt *time.Time `json:"deleteAt" `
	}

	User struct {
		UserID       int `gorm:"primary_key;unique_index"`
		FirstName    string
		LastName     string
		Email        string
		Username     string
		UserPassword string
		PhoneNo      string
		RoleID       int `gorm:"foreignkey:RoleID"`
		ProfileImage string
		AuthStatus   int
		GormDefault
	}

	Role struct {
		RoleID   int    `gorm:"primary_key;unique_index"`
		RoleName string `gorm:"type:varchar(30);unique_index"`
		GormDefault
	}

	APIToken struct {
		TokenID int    `gorm:"primary_key;unique_index"`
		UserID  string `gorm:"type:varchar(30);unique_index"`
		GormDefault
	}

	NotifyMethod struct {
		NotifyID string `gorm:"type:varchar(15);primary_key;unique_index"`
		EmpID    int    `gorm:"foreignkey:EmpID"`
		Key      string `gorm:"type:varchar(20);unique_index"`
		Value    string `gorm:"type:varchar(255);unique_index"`
		GormDefault
	}

	////// CUSTOM QUERY /////

	AuditLog struct {
		AuditID    int `gorm:"primary_key;unique_index"`
		EmpID      string
		EmpName    string
		Department string
		Argument   string
		ModuleName string
		API        string
		GormDefault
	}
)
