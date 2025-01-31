package audit

import (
	"fmt"
	"strings"

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

// ToString implements ResourceAudit.
func (p *AccountLevelAudit) ToString() string {
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
func (p *AccountLevelAudit) SubmitAuditDetail(message string) error {
	if p.auditService == nil {
		fmt.Println("auditService is nil, silent fail")
		return nil
	}

	p.details = append(p.details, message)
	return nil
}

// SubmitAuditLog implements ResourceAudit.
// This function is used to submit an audit log for an account level resource.
// Calling this function will automatically set the parentAuditLogId field.
// parentAuditLogId is used to log audit details to a parent audit log.
func (p *AccountLevelAudit) SubmitAuditLog(action auditmodels.Action) error {
	if p.auditService == nil {
		fmt.Println("auditService is nil, silent fail")
		return nil
	}

	var logId uuid.UUID
	var err error

	logId, err = p.auditService.PutAuditLogAccountLevel(p.resourceType, p.resourceId, p.metadata.Name, action, p.userId, p.ownerId, p.details)
	if err != nil {
		return err
	}

	if logId == uuid.Nil {
		return fmt.Errorf("logId is nil")
	}

	p.parentAuditLogId = logId.String()
	p.actionTaken = action
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
