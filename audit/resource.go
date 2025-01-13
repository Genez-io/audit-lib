package audit

import (
	auditmodels "github.com/Genez-io/audit-lib/audit_models"
)

type BaseProjectLevelAudit struct {
	auditService *auditService
	resourceType auditmodels.ProjectLevelResource
	userId       string
	resourceId   string
	projectId    string
	ownerId      *string
}

type BaseAccountLevelAudit struct {
	auditService *auditService
	resourceType auditmodels.AccountLevelResource
	userId       string
	resourceId   string
	ownerId      *string
}

type ResourceAudit interface {
	SubmitAuditLog(action auditmodels.Action) (string, error)
	SubmitAuditDetail(parentId string, message string) (string, error)
}
