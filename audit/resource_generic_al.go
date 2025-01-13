package audit

import (
	auditmodels "github.com/Genez-io/audit-lib/audit_models"
	"github.com/google/uuid"
)

type GenericAccountMetadata struct {
	Name string
}

type GenericAccountAudit struct {
	BaseAccountLevelAudit
	// Specific fields for this resource.
	metadata GenericAccountMetadata
}

// SubmitAuditDetail implements ResourceAudit.
func (p *GenericAccountAudit) SubmitAuditDetail(parentId string, message string) (string, error) {
	parsedParentId, err := uuid.Parse(parentId)
	if err != nil {
		return "", err
	}
	logId, err := GetAuditService().PutAuditDetail(parsedParentId, message)
	return logId.String(), err
}

// SubmitAuditLog implements ResourceAudit.
func (p *GenericAccountAudit) SubmitAuditLog(action auditmodels.Action) (string, error) {
	var logId uuid.UUID
	var err error

	logId, err = p.auditService.PutAuditLogAccountLevel(p.resourceType, p.resourceId, p.metadata.Name, action, p.userId, p.ownerId)
	return logId.String(), err
}

func NewGenericAccountAudit(resType auditmodels.AccountLevelResource, userId, resourceId, resourceName string, ownerId *string) ResourceAudit {
	return &GenericAccountAudit{
		BaseAccountLevelAudit: BaseAccountLevelAudit{
			resourceType: resType,
			userId:       userId,
			auditService: GetAuditService(),
			ownerId:      ownerId,
			resourceId:   resourceId,
		},
		metadata: GenericAccountMetadata{
			Name: resourceName,
		},
	}
}
