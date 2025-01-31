package audit

import (
	"fmt"

	auditmodels "github.com/Genez-io/audit-lib/audit_models"
	"github.com/Genez-io/audit-lib/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type auditService struct {
	repo *repository.AuditRepository
}

var auditServiceInstance *auditService = nil

func NewAuditServiceWithDB(db *gorm.DB) *auditService {
	if db == nil {
		panic("db is required")
	}

	if auditServiceInstance != nil {
		return auditServiceInstance
	}

	auditRepository := repository.NewAuditRepository(db)
	auditServiceInstance = &auditService{
		repo: auditRepository,
	}

	return auditServiceInstance
}

func GetAuditService() *auditService {
	if auditServiceInstance == nil {
		fmt.Println("auditServiceInstance is nil, all calls to audit service will silently fail")
		return nil
	}

	return auditServiceInstance
}

func (s *auditService) PutAuditLogProjectLevel(resource auditmodels.ProjectLevelResource, resourceId, resourceName string, action auditmodels.Action, userId, projectId string, details []string) (uuid.UUID, error) {
	auditLog, err := s.repo.CreateAuditLogProjectLevel(string(resource), resourceId, resourceName, string(action), userId, projectId, details)
	if err != nil {
		return uuid.UUID{}, err
	}
	return auditLog.ID, nil
}

func (s *auditService) PutAuditLogAccountLevel(resourceType auditmodels.AccountLevelResource, resourceId, resourceName string, action auditmodels.Action, userId string, orgId *string, details []string) (uuid.UUID, error) {
	auditLog, err := s.repo.CreateAuditLogAccountLevel(string(resourceType), resourceId, resourceName, string(action), userId, orgId, details)
	if err != nil {
		return uuid.UUID{}, err
	}

	return auditLog.ID, nil
}
