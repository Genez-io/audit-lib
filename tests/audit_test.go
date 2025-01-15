package tests

import (
	"fmt"
	"testing"

	"github.com/Genez-io/audit-lib/audit"
	auditmodels "github.com/Genez-io/audit-lib/audit_models"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAuditRepository_PutGenericAccountAuditLog(t *testing.T) {
	dbuser := "genezio"
	dbpass := "genezio"
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/genezio?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpass)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	audit.NewAuditServiceWithDB(db)
	userId := uuid.New().String()
	resourceId := uuid.New().String()
	pa := audit.NewGenericAccountAudit(auditmodels.Users, userId, resourceId, "test-generic-al", nil)
	logId, err := pa.SubmitAuditLog(auditmodels.ActionCreate)
	if err != nil {
		t.Error(err)
	}

	logDetailId, err := pa.SubmitAuditDetail(logId, "test message")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(logId, logDetailId)
}

func TestAuditRepository_PutGenericProjectAuditLog(t *testing.T) {
	dbuser := "genezio"
	dbpass := "genezio"
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/genezio?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpass)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	audit.NewAuditServiceWithDB(db)

	userId := uuid.New().String()
	projectId := uuid.New().String()
	resourceId := uuid.New().String()
	pa := audit.NewGenericProjectAudit(auditmodels.Deployments, userId, projectId, resourceId, "test-generic-pl", nil)
	logId, err := pa.SubmitAuditLog(auditmodels.ActionCreate)
	if err != nil {
		t.Error(err)
	}

	logDetailId, err := pa.SubmitAuditDetail(logId, "test message")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(logId, logDetailId)
}
