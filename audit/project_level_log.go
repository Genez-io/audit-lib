package audit

import (
	"fmt"
	"strings"

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

// ToString implements ResourceAudit.
func (p *ProjectLevelAudit) ToString() string {
	var message string

	singularResource := p.resourceType[:len(p.resourceType)-1]
	singularCamelCaseResource := strings.ToUpper(string(p.resourceType[:1])) + string(p.resourceType[1:len(p.resourceType)-1])

	switch p.actionTaken {
	case auditmodels.ActionCreate:
		if p.metadata.Name == "" {
			message = fmt.Sprintf("Created a new %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was created", singularCamelCaseResource, p.metadata.Name)
		}
	case auditmodels.ActionUpdate:
		if p.metadata.Name == "" {
			message = fmt.Sprintf("Updated a %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was updated", singularCamelCaseResource, p.metadata.Name)
		}
	case auditmodels.ActionDelete:
		if p.metadata.Name == "" {
			message = fmt.Sprintf("Deleted a %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was deleted", singularCamelCaseResource, p.metadata.Name)
		}
	case auditmodels.ActionDisable:
		if p.metadata.Name == "" {
			message = fmt.Sprintf("Disabled a %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was disabled", singularCamelCaseResource, p.metadata.Name)
		}
	case auditmodels.ActionEnable:
		if p.metadata.Name == "" {
			message = fmt.Sprintf("Enabled a %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was enabled", singularCamelCaseResource, p.metadata.Name)
		}
	case auditmodels.ActionLogin:
		return "User logged in"
	case auditmodels.ActionLogout:
		return "User logged out"
	default:
		message = fmt.Sprintf("Invalid action: %s", p.actionTaken)
	}

	return message
}

// SubmitAuditDetail implements ResourceAudit.
func (p *ProjectLevelAudit) SubmitAuditDetail(message string) error {
	if p.auditService == nil {
		fmt.Println("auditService is nil, silent fail")
		return nil
	}
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
	if p.auditService == nil {
		fmt.Println("auditService is nil, silent fail")
		return nil
	}
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
	p.actionTaken = action
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
