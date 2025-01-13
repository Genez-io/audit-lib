package audit

import (
	"fmt"

	auditmodels "github.com/Genez-io/audit-lib/audit_models"
	"github.com/google/uuid"
)

type AccountLevelMetadata struct {
	Name string
}

type AccountLevelAudit struct {
	BaseAccountLevelAudit
	// Specific fields for this resource.
	metadata AccountLevelMetadata
}

// SubmitAuditDetail implements ResourceAudit.
func (p *AccountLevelAudit) SubmitAuditDetail(message string) error {
	if p.parentAuditLogId == "" {
		return fmt.Errorf("parentAuditLogId is not set")
	}
	parsedParentId, err := uuid.Parse(p.parentAuditLogId)
	if err != nil {
		return err
	}
	_, err = GetAuditService().PutAuditDetail(parsedParentId, message)
	return err
}

// SubmitAuditLog implements ResourceAudit.
// This function is used to submit an audit log for an account level resource.
// Calling this function will automatically set the parentAuditLogId field.
// parentAuditLogId is used to log audit details to a parent audit log.
func (p *AccountLevelAudit) SubmitAuditLog(action auditmodels.Action) error {
	var logId uuid.UUID
	var err error

	logId, err = p.auditService.PutAuditLogAccountLevel(p.resourceType, p.resourceId, p.metadata.Name, action, p.userId, p.ownerId)
	if err != nil {
		return err
	}

	if logId == uuid.Nil {
		return fmt.Errorf("logId is nil")
	}

	p.parentAuditLogId = logId.String()
	return nil
}

func NewAccountLevelAudit(resType auditmodels.AccountLevelResource, userId, resourceId, resourceName string, ownerId *string) ResourceAudit {
	return &AccountLevelAudit{
		BaseAccountLevelAudit: BaseAccountLevelAudit{
			resourceType: resType,
			userId:       userId,
			auditService: GetAuditService(),
			ownerId:      ownerId,
			resourceId:   resourceId,
		},
		metadata: AccountLevelMetadata{
			Name: resourceName,
		},
	}
}
