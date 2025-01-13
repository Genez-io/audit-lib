package audit

import (
	"fmt"

	auditmodels "github.com/Genez-io/audit-lib/audit_models"
	"github.com/google/uuid"
)

type ProjectLevelMetadata struct {
	Name string
}

type ProjectLevelAudit struct {
	BaseProjectLevelAudit
	// Specific fields for this resource.
	metadata ProjectLevelMetadata
}

// SubmitAuditDetail implements ResourceAudit.
func (p *ProjectLevelAudit) SubmitAuditDetail(message string) error {
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
// This function is used to submit an audit log for a project level resource.
// Calling this function will automatically set the parentAuditLogId field.
// parentAuditLogId is used to log audit details to a parent audit log.
func (p *ProjectLevelAudit) SubmitAuditLog(action auditmodels.Action) error {
	var logId uuid.UUID
	var err error
	// Check that the action is valid.
	logId, err = p.auditService.PutAuditLogProjectLevel(p.resourceType, p.resourceId, p.metadata.Name, action, p.userId, p.projectId)
	if err != nil {
		return err
	}

	if logId == uuid.Nil {
		return fmt.Errorf("logId is nil")
	}

	p.parentAuditLogId = logId.String()
	return err
}

func NewProjectLevelAudit(resType auditmodels.ProjectLevelResource, userId, projectId, resourceId, resourceName string, ownerId *string) ResourceAudit {
	return &ProjectLevelAudit{
		BaseProjectLevelAudit: BaseProjectLevelAudit{
			resourceType: resType,
			userId:       userId,
			auditService: GetAuditService(),
			resourceId:   resourceId,
			ownerId:      ownerId,
			projectId:    projectId,
		},
		metadata: ProjectLevelMetadata{
			Name: resourceName,
		},
	}
}
