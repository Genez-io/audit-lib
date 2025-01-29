package auditmodels

import (
	"fmt"
	"time"
)

type AccountLevelResource string

const (
	Projects      AccountLevelResource = "projects"
	Databases     AccountLevelResource = "databases"
	Billing       AccountLevelResource = "billings"
	Users         AccountLevelResource = "users"
	AccessTokens  AccountLevelResource = "access_tokens"
	Collaboration AccountLevelResource = "collaborations"
)

func (AccountLevelResource) Check(s string) error {
	switch s {
	case "projects", "databases", "billings", "users", "access_tokens", "collaborations":
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
	Envs                  ProjectLevelResource = "envs"
	DatabaseAssignments   ProjectLevelResource = "database_links"
	ClassPauses           ProjectLevelResource = "class_pauses"
	FunctionPauses        ProjectLevelResource = "function_pauses"
	CustomDomains         ProjectLevelResource = "custom_domains"
	Integrations          ProjectLevelResource = "integrations"
	Authentication        ProjectLevelResource = "authentications"
	AuthenticationMethods ProjectLevelResource = "authentication_methods"
	AuthenticationMail    ProjectLevelResource = "authentication_mails"
	EmailService          ProjectLevelResource = "email_services"
	LogDrains             ProjectLevelResource = "log_drains"
)

func (ProjectLevelResource) Check(s string) error {
	switch s {
	case
		"deployments",
		"code_updates",
		"collaborators",
		"envs",
		"database_links",
		"class_pauses",
		"function_pauses",
		"custom_domains",
		"integrations",
		"authentications",
		"authentication_methods",
		"authentication_mails",
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
	ActionLogin   Action = "login"
	ActionLogout  Action = "logout"
)

type AuditFilter struct {
	Resource     AccountLevelResource
	Action       Action
	Before       time.Time
	After        time.Time
	AuthorUserID string
}
