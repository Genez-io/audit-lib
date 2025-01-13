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
		panic("audit service not initialized")
	}

	return auditServiceInstance
}

func (s *auditService) PutAuditLogProjectLevel(resource auditmodels.ProjectLevelResource, resourceId, resourceName string, action auditmodels.Action, userId, projectId string) (uuid.UUID, error) {
	auditLog, err := s.repo.CreateAuditLogProjectLevel(string(resource), resourceId, resourceName, string(action), userId, projectId)
	if err != nil {
		return uuid.UUID{}, err
	}

	fmt.Println(auditLog)
	return auditLog.ID, nil
}
func (s *auditService) PutAuditLogAccountLevel(resourceType auditmodels.AccountLevelResource, resourceId, resourceName string, action auditmodels.Action, userId string, orgId *string) (uuid.UUID, error) {
	auditLog, err := s.repo.CreateAuditLogAccountLevel(string(resourceType), resourceId, resourceName, string(action), userId, orgId)
	if err != nil {
		return uuid.UUID{}, err
	}

	fmt.Println(auditLog)
	return auditLog.ID, nil
}

func (s *auditService) PutAuditDetail(parentId uuid.UUID, message string) (uuid.UUID, error) {
	auditLogDetail, err := s.repo.CreateAuditLogDetail(parentId, message)
	if err != nil {
		return uuid.UUID{}, err
	}

	fmt.Println(auditLogDetail)
	return auditLogDetail.ID, nil

}
