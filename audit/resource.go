package audit

import (
	"fmt"
	"strings"

	auditmodels "github.com/Genez-io/audit-lib/audit_models"
)

type BaseProjectLevelAudit struct {
	auditService *auditService
	resourceType auditmodels.ProjectLevelResource
	userId       string
	resourceId   string
	projectId    string
	ownerId      *string
	// Note this should be populated when calling SubmitAuditLog.
	parentAuditLogId string
	actionTaken      auditmodels.Action
	details          []string
}

type BaseAccountLevelAudit struct {
	auditService *auditService
	resourceType auditmodels.AccountLevelResource
	userId       string
	resourceId   string
	ownerId      *string
	// Note this should be populated when calling SubmitAuditLog.
	parentAuditLogId string
	actionTaken      auditmodels.Action
	details          []string
}

type ResourceAudit interface {
	SubmitAuditLog(action auditmodels.Action) error
	SubmitAuditDetail(message string) error
	ToString() string
}

func AuditLogToString(resourceType, resourceName, action string) string {
	var message string

	singularResource := resourceType[:len(resourceType)-1]
	singularCamelCaseResource := strings.ToUpper(string(resourceType[:1])) + string(resourceType[1:len(resourceType)-1])
	switch action {
	case string(auditmodels.ActionCreate):
		if resourceName == "" {
			message = fmt.Sprintf("Created a new %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was created", singularCamelCaseResource, resourceName)
		}
	case string(auditmodels.ActionUpdate):
		if resourceName == "" {
			message = fmt.Sprintf("Updated a %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was updated", singularCamelCaseResource, resourceName)
		}
	case string(auditmodels.ActionDelete):
		if resourceName == "" {
			message = fmt.Sprintf("Deleted a %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was deleted", singularCamelCaseResource, resourceName)
		}
	case string(auditmodels.ActionDisable):
		if resourceName == "" {
			message = fmt.Sprintf("Disabled a %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was disabled", singularCamelCaseResource, resourceName)
		}
	case string(auditmodels.ActionEnable):
		if resourceName == "" {
			message = fmt.Sprintf("Enabled a %s", singularResource)
		} else {
			message = fmt.Sprintf("%s %s was enabled", singularCamelCaseResource, resourceName)
		}
	case string(auditmodels.ActionLogin):
		return "User logged in"
	case string(auditmodels.ActionLogout):
		return "User logged out"
	default:
		message = fmt.Sprintf("Invalid action: %s", action)
	}

	return message
}
