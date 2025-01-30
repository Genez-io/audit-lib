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
}

type ResourceAudit interface {
	SubmitAuditLog(action auditmodels.Action) error
	SubmitAuditDetail(message string) error
	ToString() string
}

func AuditLogToString(resourceType, resourceName, action string) string {
	var message string

	resourceTypeParsed := strings.ReplaceAll(resourceType, "_", " ")

	singularResource := resourceTypeParsed[:len(resourceTypeParsed)-1]
	singularCamelCaseResource := strings.ToUpper(string(resourceTypeParsed[:1])) + string(resourceTypeParsed[1:len(resourceTypeParsed)-1])
	switch action {
	case string(auditmodels.ActionCreate):
		switch resourceType {
		case string(auditmodels.Envs):
			message = fmt.Sprintf("Environment Variable %s was created", resourceName)
		case string(auditmodels.DatabaseAssignments):
			message = fmt.Sprintf("Database %s was linked", resourceName)
		case string(auditmodels.Collaborators):
			message = fmt.Sprintf("Collaborator %s was invited to the project", resourceName)
		default:
			if resourceName == "" {
				message = fmt.Sprintf("Created a new %s", singularResource)
			} else {
				message = fmt.Sprintf("%s %s was created", singularCamelCaseResource, resourceName)
			}
		}
	case string(auditmodels.ActionUpdate):
		switch resourceType {
		case string(auditmodels.Envs):
			message = fmt.Sprintf("Environment Variable %s was modified", resourceName)
		case string(auditmodels.Deployments):
			message = fmt.Sprintf("Deployed project %s", resourceName)
		case string(auditmodels.CodeUpdates):
			message = fmt.Sprintf("Code updated for project %s", resourceName)
		default:
			if resourceName == "" {
				message = fmt.Sprintf("Updated %s", singularResource)
			} else {
				message = fmt.Sprintf("%s %s was updated", singularCamelCaseResource, resourceName)
			}
		}
	case string(auditmodels.ActionDelete):
		switch resourceType {
		case string(auditmodels.Envs):
			message = fmt.Sprintf("Environment Variable %s was deleted", resourceName)
		default:
			if resourceName == "" {
				message = fmt.Sprintf("Deleted %s", singularResource)
			} else {
				message = fmt.Sprintf("%s %s was deleted", singularCamelCaseResource, resourceName)
			}
		}
	case string(auditmodels.ActionDisable):
		switch resourceType {
		default:
			messagePrefix := "Function"
			if resourceTypeParsed == "Class pause" {
				messagePrefix = "Class"
			}
			if resourceName == "" {
				message = fmt.Sprintf("Disabled %s", singularResource)
			} else {
				message = fmt.Sprintf("%s \"%s\" was disabled", messagePrefix, resourceName)
			}
		}
	case string(auditmodels.ActionEnable):
		switch resourceType {
		default:
			messagePrefix := "Function"
			if resourceTypeParsed == "Class pause" {
				messagePrefix = "Class"
			}
			if resourceName == "" {
				message = fmt.Sprintf("Enabled %s", singularResource)
			} else {
				message = fmt.Sprintf("%s \"%s\" was enabled", messagePrefix, resourceName)
			}
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
