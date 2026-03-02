package services

import (
	"strings"
	"testing"

	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/internal/users/forms"
	"github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUsersServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE insight_orgs (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			parent_id INTEGER NOT NULL,
			key TEXT,
			level INTEGER NOT NULL,
			path TEXT,
			description TEXT,
			creator TEXT,
			updater TEXT,
			created_at DATETIME,
			updated_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("create insight_orgs table failed: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE insight_org_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uid INTEGER NOT NULL,
			organization_key TEXT NOT NULL,
			role_id INTEGER,
			is_active BOOLEAN,
			created_at DATETIME,
			updated_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("create insight_org_users table failed: %v", err)
	}

	oldDB := global.App.DB
	oldLog := global.App.Log
	global.App.DB = db
	global.App.Log = logrus.New()
	t.Cleanup(func() {
		global.App.DB = oldDB
		global.App.Log = oldLog
	})

	return db
}

func TestApplyOrgDescendantScope(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}

	var rows []models.InsightOrgUsers
	statement := applyOrgDescendantScope(db.Table("insight_org_users b"), "b.organization_key", "0-1").Find(&rows).Statement

	sql := statement.SQL.String()
	if !strings.Contains(sql, "b.organization_key = ? OR b.organization_key LIKE ?") {
		t.Fatalf("unexpected sql condition: %s", sql)
	}
	if len(statement.Vars) != 2 {
		t.Fatalf("unexpected vars length: %d", len(statement.Vars))
	}
	if statement.Vars[0] != "0-1" || statement.Vars[1] != "0-1-%" {
		t.Fatalf("unexpected vars: %#v", statement.Vars)
	}
}

func TestDeleteOrganizationsServiceRunDeletesDescendantsAndBindings(t *testing.T) {
	db := setupUsersServiceTestDB(t)

	root := models.InsightOrgs{ID: 1, Name: "基础架构部", ParentID: 0, Key: "0-1", Level: 1}
	child := models.InsightOrgs{
		ID:       2,
		Name:     "数据库组",
		ParentID: 1,
		Key:      "0-1-2",
		Level:    2,
		Path:     datatypes.JSON([]byte(`["0-1"]`)),
	}
	grandChild := models.InsightOrgs{
		ID:       3,
		Name:     "MySQL小组",
		ParentID: 2,
		Key:      "0-1-2-3",
		Level:    3,
		Path:     datatypes.JSON([]byte(`["0-1","0-1-2"]`)),
	}
	otherRoot := models.InsightOrgs{ID: 10, Name: "质量保障部", ParentID: 0, Key: "0-10", Level: 1}

	if err := db.Create(&[]models.InsightOrgs{root, child, grandChild, otherRoot}).Error; err != nil {
		t.Fatalf("create organizations failed: %v", err)
	}

	bindings := []models.InsightOrgUsers{
		{Uid: 1001, OrganizationKey: "0-1", RoleID: 1},
		{Uid: 1002, OrganizationKey: "0-1-2-3", RoleID: 2},
		{Uid: 1003, OrganizationKey: "0-10", RoleID: 3},
	}
	if err := db.Create(&bindings).Error; err != nil {
		t.Fatalf("create bindings failed: %v", err)
	}

	service := DeleteOrganizationsService{
		DeleteOrganizationsForm: &forms.DeleteOrganizationsForm{
			Key:  "0-1",
			Name: "基础架构部",
		},
	}
	if err := service.Run(); err != nil {
		t.Fatalf("delete organization failed: %v", err)
	}

	var orgs []models.InsightOrgs
	if err := db.Order("id asc").Find(&orgs).Error; err != nil {
		t.Fatalf("query organizations failed: %v", err)
	}
	if len(orgs) != 1 || orgs[0].Key != "0-10" {
		t.Fatalf("unexpected organizations after delete: %#v", orgs)
	}

	var remainingBindings []models.InsightOrgUsers
	if err := db.Order("uid asc").Find(&remainingBindings).Error; err != nil {
		t.Fatalf("query bindings failed: %v", err)
	}
	if len(remainingBindings) != 1 || remainingBindings[0].OrganizationKey != "0-10" {
		t.Fatalf("unexpected bindings after delete: %#v", remainingBindings)
	}
}
