package model

import (
	"time"

	"gorm.io/gorm"
)

type RoleType string

const (
	Director          RoleType = "director"           // Director
	TradeTeam         RoleType = "trade-team"         // Trade Team
	NSM               RoleType = "nsm"                // National Sales Manager
	SM                RoleType = "sm"                 // Sales Manager
	ASM               RoleType = "asm"                // Assistant Sales Manager
	FieldForce        RoleType = "field-force"        // Field Force
	FIC               RoleType = "fic"                // Finance Internal Control
	MarketingDirector RoleType = "marketing-director" // Marketing Director
	MSD               RoleType = "msd"                // Marketing Service Department
)

type Roles struct {
	ID        *uint          `json:"role_id" gorm:"primaryKey"`
	Name      *string        `json:"role_name" gorm:"type:varchar(256); not null"`
	Label     *string        `json:"role_label" gorm:"type:varchar(256); not null"`
	ParentId  *uint          `json:"parent_id"`
	CreatedAt time.Time      `gorm:"DEFAULT:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type ParentRoles struct {
	ID        *uint      `json:"role_id" gorm:"primaryKey"`
	Name      *string    `json:"role_name" gorm:"type:varchar(256); not null"`
	Label     *string    `json:"role_label" gorm:"type:varchar(256); not null"`
	CreatedAt *time.Time `gorm:"DEFAULT:current_timestamp" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (ParentRoles) TableName() string {
	return "roles"
}

type RolesRawQuerySelfJoinResult struct {
	RoleId       *uint
	Name         *string
	Label        *string
	ParentId     *uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
	NameRoleHead *string
}

type RolesWithEpochEntity struct {
	RoleId    *uint   `json:"role_id"`
	Name      *string `json:"name"`
	Label     *string `json:"label"`
	ParentId  *uint   `json:"parent_id"`
	CreatedAt *uint64 `json:"created_at"`
	UpdatedAt *uint64 `json:"updated_at"`
}
