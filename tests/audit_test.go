package tests

import (
	"fmt"
	"testing"
	"time"

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
	pa := audit.NewAccountLevelAudit(auditmodels.Users, userId, resourceId, "test-generic-al", nil)
	if err != nil {
		t.Error(err)
	}
	err = pa.SubmitAuditDetail("test message1")
	if err != nil {
		t.Error(err)
	}

	pa.SubmitAuditDetail("test message2")

	err = pa.SubmitAuditLog(auditmodels.ActionCreate)
	if err != nil {
		t.Error(err)
	}

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
	pa := audit.NewProjectLevelAudit(auditmodels.Deployments, userId, projectId, resourceId, "test-generic-pl", nil)
	err = pa.SubmitAuditLog(auditmodels.ActionCreate)
	if err != nil {
		t.Error(err)
	}

	err = pa.SubmitAuditDetail("test message")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pa.ToString())
}

func TestAuditRepository_AuditLogRateLimited(t *testing.T) {
	dbuser := "genezio"
	dbpass := "genezio"
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/genezio?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpass)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	audit.NewAuditServiceWithDB(db)

	for range 100 {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Pushing audit log")
		userId := uuid.New().String()
		resourceId := uuid.New().String()
		pa := audit.NewAccountLevelAudit(auditmodels.Users, userId, resourceId, "test-generic-al", nil)
		if err != nil {
			t.Error(err)
		}
		err = pa.SubmitAuditDetail("test message1")
		if err != nil {
			t.Error(err)
		}

		pa.SubmitAuditDetail("test message2")

		err = pa.SubmitAuditLog(auditmodels.ActionCreate)
		if err != nil {
			t.Error(err)
		}
	}

	time.Sleep(10 * time.Minute)

}
