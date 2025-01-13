package auditmodels

import (
	"fmt"
	"time"
)

type AccountLevelResource string

const (
	Projects      AccountLevelResource = "projects"
	Databases     AccountLevelResource = "databases"
	Billing       AccountLevelResource = "billing"
	Users         AccountLevelResource = "users"
	AccessTokens  AccountLevelResource = "access_tokens"
	Collaboration AccountLevelResource = "collaborations"
)

func (AccountLevelResource) Check(s string) error {
	switch s {
	case "projects", "databases", "billing", "users", "access_tokens", "collaborations":
		return nil
	default:
		return fmt.Errorf("invalid AccountLevelResource: %s", s)
	}
}

type ProjectLevelResource string

const (
	Deployments           ProjectLevelResource = "deployments"
	CodeUpdates           ProjectLevelResource = "code_updates"
	Collaborators         ProjectLevelResource = "collaborators"
	Environments          ProjectLevelResource = "environments"
	DatabaseAssignments   ProjectLevelResource = "databases"
	ClassPauses           ProjectLevelResource = "class_pauses"
	CustomDomains         ProjectLevelResource = "custom_domains"
	Integrations          ProjectLevelResource = "integrations"
	Authentication        ProjectLevelResource = "authentication"
	AuthenticationMethods ProjectLevelResource = "authentication_methods"
	AuthenticationMail    ProjectLevelResource = "authentication_mail"
	EmailService          ProjectLevelResource = "email_service"
	LogDrains             ProjectLevelResource = "log_drains"
)

func (ProjectLevelResource) Check(s string) error {
	switch s {
	case
		"deployments",
		"code_updates",
		"collaborators",
		"environments",
		"databases",
		"class_pauses",
		"custom_domains",
		"integrations",
		"authentication",
		"authentication_methods",
		"authentication_mail",
		"email_service",
		"log_drains":
		return nil
	default:
		return fmt.Errorf("invalid ProjectLevelResource: %s", s)
	}
}

type Action string

const (
	ActionCreate  Action = "create"
	ActionRead    Action = "read"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionEnable  Action = "enable"
	ActionDisable Action = "disable"
)

type AuditFilter struct {
	Resource     AccountLevelResource
	Action       Action
	Before       time.Time
	After        time.Time
	AuthorUserID string
}
