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
	// Note this should be populated when calling SubmitAuditLog.
	parentAuditLogId string
	actionTaken 		auditmodels.Action
}

type BaseAccountLevelAudit struct {
	auditService *auditService
	resourceType auditmodels.AccountLevelResource
	userId       string
	resourceId   string
	ownerId      *string
	// Note this should be populated when calling SubmitAuditLog.
	parentAuditLogId string
	actionTaken 		auditmodels.Action
}

type ResourceAudit interface {
	SubmitAuditLog(action auditmodels.Action) error
	SubmitAuditDetail(message string) error
	ToString() string
}
