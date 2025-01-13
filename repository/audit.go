package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditRepository struct {
	DB *gorm.DB
}

type AuditLevel string

const (
	AuditLevelAccount AuditLevel = "account"
	AuditLevelProject AuditLevel = "project"
)

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:char(36);column:id;primaryKey;default:(UUID())"`
	CreatedAt  time.Time
	Level      AuditLevel `gorm:"type:enum('account', 'project')"`
	ResourceID string     `gorm:"index"`
	Name       string     `gorm:"index"`

	// User who performed the action
	UserID string `gorm:"index"`

	// Composite index for resource_type:action pairs (ex: create:project)
	Action       string `gorm:"index:idx_action_resource"`
	ResourceType string `gorm:"index:idx_action_resource"`

	// If the resource is account level, this will be the owner's ID.
	// If the owner id is nil, then the owner is the user, otherwise the owner is the organization.
	OwnerID string `gorm:"index"`

	// If the resource is project level, this will be the project's ID.
	ProjectID string `gorm:"index"`
}

type AuditLogDetail struct {
	ID         uuid.UUID `gorm:"type:char(36);column:id;primaryKey;default:(UUID())"`
	CreatedAt  time.Time
	AuditLogID uuid.UUID `gorm:"index"`
	AuditLog   AuditLog  `gorm:"constraint:OnDelete:CASCADE"`
	Message    string
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	db.AutoMigrate(&AuditLog{}, &AuditLogDetail{})
	return &AuditRepository{DB: db}
}

func (r *AuditRepository) CreateAuditLogProjectLevel(
	resourceType string,
	resourceId string,
	resourceName string,
	action string,
	userId string,
	projectId string,
) (AuditLog, error) {
	auditLog := AuditLog{
		ID:           uuid.New(),
		ResourceType: resourceType,
		ResourceID:   resourceId,
		Action:       action,
		UserID:       userId,
		ProjectID:    projectId,
		Name:         resourceName,
		Level:        AuditLevelProject,
	}

	res := r.DB.Create(&auditLog)
	return auditLog, res.Error
}

func (r *AuditRepository) CreateAuditLogAccountLevel(
	resourceType string,
	resourceId string,
	resourceName string,
	action string,
	userId string,
	ownerId *string,
) (AuditLog, error) {
	auditLog := AuditLog{
		ID:           uuid.New(),
		ResourceType: resourceType,
		ResourceID:   resourceId,
		Name:         resourceName,
		Action:       action,
		UserID:       userId,
		Level:        AuditLevelAccount,
	}

	if ownerId != nil {
		auditLog.OwnerID = *ownerId
	}

	res := r.DB.Create(&auditLog)

	return auditLog, res.Error
}

func (r *AuditRepository) CreateAuditLogDetail(
	parentId uuid.UUID,
	message string,
) (AuditLogDetail, error) {
	auditLogDetail := AuditLogDetail{
		ID:         uuid.New(),
		AuditLogID: parentId,
		Message:    message,
	}

	err := r.DB.Create(&auditLogDetail).Error
	return auditLogDetail, err
}
