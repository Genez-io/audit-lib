package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

type AuditRepository struct {
	db            *gorm.DB
	auditLogQueue []AuditLog
	queueLock     sync.Mutex
	rateLimiter   *rate.Limiter
	context       context.Context
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

	Details []AuditLogDetail
}

type AuditLogDetail struct {
	ID         uuid.UUID `gorm:"type:char(36);column:id;primaryKey;default:(UUID())"`
	CreatedAt  time.Time
	AuditLogID uuid.UUID `gorm:"index"`
	AuditLog   AuditLog  `gorm:"constraint:OnDelete:CASCADE"`
	Message    string
}

const batchSize = 50
const flushInterval = 5 * time.Second

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	tickerFlush := time.NewTicker(flushInterval)
	done := make(chan bool)

	auditRepository := &AuditRepository{
		db:            db,
		auditLogQueue: make([]AuditLog, 0),
		queueLock:     sync.Mutex{},
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-tickerFlush.C:
				if len(auditRepository.auditLogQueue) > 0 {
					auditRepository.FlushAuditLogQueue()
				}
			}
		}
	}()
	return auditRepository
}

func (r *AuditRepository) CreateAuditLogProjectLevel(
	resourceType string,
	resourceId string,
	resourceName string,
	action string,
	userId string,
	projectId string,
	details []string,
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

	for _, detail := range details {
		auditLog.Details = append(auditLog.Details, AuditLogDetail{
			ID:      uuid.New(),
			Message: detail,
		})
	}
	r.PushToAuditQueue(auditLog)
	return auditLog, nil
}

func (r *AuditRepository) CreateAuditLogAccountLevel(
	resourceType string,
	resourceId string,
	resourceName string,
	action string,
	userId string,
	ownerId *string,
	details []string,
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

	for _, detail := range details {
		auditLog.Details = append(auditLog.Details, AuditLogDetail{
			ID:      uuid.New(),
			Message: detail,
		})
	}

	r.PushToAuditQueue(auditLog)
	return auditLog, nil
}

func (r *AuditRepository) PushToAuditQueue(auditLog AuditLog) {
	r.queueLock.Lock()
	r.auditLogQueue = append(r.auditLogQueue, auditLog)
	r.queueLock.Unlock()
}

func (r *AuditRepository) FlushAuditLogQueue() error {
	// Create a deep copy of the previous queue.
	r.queueLock.Lock()
	var oldQueue []AuditLog = make([]AuditLog, 0)
	for _, log := range r.auditLogQueue {
		var copyLog AuditLog = log
		copy(copyLog.Details, log.Details)
		oldQueue = append(oldQueue, copyLog)
	}
	r.auditLogQueue = make([]AuditLog, 0)
	r.queueLock.Unlock()

	res := r.db.CreateInBatches(&oldQueue, batchSize)
	return res.Error
}
