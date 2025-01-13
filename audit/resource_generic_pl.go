package audit

import (
	auditmodels "github.com/Genez-io/audit-lib/audit_models"
	"github.com/google/uuid"
)

type GenericProjectMetadata struct {
	Name string
}

type GenericProjectAudit struct {
	BaseProjectLevelAudit
	// Specific fields for this resource.
	metadata GenericProjectMetadata
}

// SubmitAuditDetail implements ResourceAudit.
func (p *GenericProjectAudit) SubmitAuditDetail(parentId string, message string) (string, error) {
	parsedParentId, err := uuid.Parse(parentId)
	if err != nil {
		return "", err
	}
	logId, err := GetAuditService().PutAuditDetail(parsedParentId, message)
	return logId.String(), err
}

// SubmitAuditLog implements ResourceAudit.
func (p *GenericProjectAudit) SubmitAuditLog(action auditmodels.Action) (string, error) {
	var logId uuid.UUID
	var err error
	// Check that the action is valid.
	logId, err = p.auditService.PutAuditLogProjectLevel(p.resourceType, p.resourceId, p.metadata.Name, action, p.userId, p.projectId)
	return logId.String(), err
}

func NewGenericProjectAudit(resType auditmodels.ProjectLevelResource, userId, projectId, resourceId, resourceName string, ownerId *string) ResourceAudit {
	return &GenericProjectAudit{
		BaseProjectLevelAudit: BaseProjectLevelAudit{
			resourceType: resType,
			userId:       userId,
			auditService: GetAuditService(),
			resourceId:   resourceId,
			ownerId:      ownerId,
			projectId:    projectId,
		},
		metadata: GenericProjectMetadata{
			Name: resourceName,
		},
	}
}
