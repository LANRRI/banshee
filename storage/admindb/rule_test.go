// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
)

func TestRules(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Add proj.
	proj := &models.Project{Name: "banshee"}
	assert.Ok(t, nil == db.AddProject(proj))
	// Add rules.
	rule1 := &models.Rule{Pattern: "counter.*"}
	rule2 := &models.Rule{Pattern: "timer.*"}
	assert.Ok(t, nil == db.AddRuleToProject(proj, rule1))
	assert.Ok(t, nil == db.AddRuleToProject(proj, rule2))
	// Get rules.
	rules := db.Rules()
	assert.Ok(t, len(rules) == 2)
	r1, r2 := rules[0], rules[1]
	assert.Ok(t, (r1.Equal(rule1) && r2.Equal(rule2)) || (r1.Equal(rule2) && r2.Equal(rule1)))
}

func TestGetRule(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Add proj.
	proj := &models.Project{Name: "banshee"}
	assert.Ok(t, nil == db.AddProject(proj))
	// Add rule.
	rule := &models.Rule{Pattern: "c.*"}
	assert.Ok(t, nil == db.AddRuleToProject(proj, rule))
	assert.Ok(t, rule.ID >= 1)
	// Get rule.
	r := &models.Rule{ID: rule.ID}
	assert.Ok(t, nil == db.GetRule(r))
	assert.Ok(t, r.Equal(rule))
}

func TestAddRuleToProject(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Add proj.
	proj := &models.Project{Name: "b"}
	assert.Ok(t, nil == db.AddProject(proj))
	// Add rule.
	rule := &models.Rule{Pattern: "t.*"}
	assert.Ok(t, nil == db.AddRuleToProject(proj, rule))
	// Must proj has the rule now.
	assert.Ok(t, len(proj.Rules) == 1)
	assert.Ok(t, proj.Rules[0].Equal(rule))
	// Must rule is in cache.
	r := &models.Rule{ID: rule.ID}
	assert.Ok(t, nil == db.GetRule(r))
	assert.Ok(t, r.Equal(rule))
	// Must rule in db and its project id is proj.ID
	r1 := &models.Rule{}
	err := db.db.Find(r1, rule.ID).Error
	assert.Ok(t, err == nil)
	assert.Ok(t, r1.ProjectID == proj.ID)
}

func TestDeleteRule(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Add proj.
	proj := &models.Project{Name: "b"}
	assert.Ok(t, nil == db.AddProject(proj))
	// Add rule.
	rule := &models.Rule{Pattern: "g.*"}
	assert.Ok(t, nil == db.AddRuleToProject(proj, rule))
	// Delete it.
	assert.Ok(t, nil == db.DeleteRule(rule.ID))
	// Must rule not in cache.
	assert.Ok(t, !db.HasRule(rule.ID))
	// Must rule not in db.
	err := db.db.Find(&models.Rule{}, rule.ID).Error
	assert.Ok(t, err == gorm.RecordNotFound)
}
