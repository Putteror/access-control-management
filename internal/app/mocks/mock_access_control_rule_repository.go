package mocks

import (
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"github.com/stretchr/testify/mock"
)

// MockAccessControlRuleRepository คือโครงสร้าง Mock ที่จำลองการทำงานของ AccessControlRuleRepository Interface
type MockAccessControlRuleRepository struct {
	mock.Mock
}

// PaginatedSearch คือการจำลอง method PaginatedSearch
func (m *MockAccessControlRuleRepository) GetAll(searchQuery schema.AccessControlRuleSearchQuery) ([]model.AccessControlRule, error) {
	// บอกให้ mock บันทึกการเรียก method นี้
	args := m.Called(searchQuery)

	// คืนค่าตามที่ถูกตั้งค่าไว้ (args.Get(0) คือ []model.AccessControlRule, args.Get(1) คือ int64, args.Error(2) คือ error)
	// ต้องตรวจสอบ nil ก่อนแปลงค่า
	var rules []model.AccessControlRule
	if args.Get(0) != nil {
		rules = args.Get(0).([]model.AccessControlRule)
	}

	return rules, args.Error(2)
}

// FindByID คือการจำลอง method FindByID
func (m *MockAccessControlRuleRepository) GetByID(id string) (*model.AccessControlRule, error) {
	args := m.Called(id)

	// คืนค่า rule และ error
	var rule *model.AccessControlRule
	if args.Get(0) != nil {
		rule = args.Get(0).(*model.AccessControlRule)
	}
	return rule, args.Error(1)
}

// Create คือการจำลอง method Create
func (m *MockAccessControlRuleRepository) Create(rule *model.AccessControlRule) error {
	args := m.Called(rule)
	return args.Error(0)
}

// Update คือการจำลอง method Update
func (m *MockAccessControlRuleRepository) Update(rule *model.AccessControlRule) error {
	args := m.Called(rule)
	return args.Error(0)
}

// Delete คือการจำลอง method Delete
func (m *MockAccessControlRuleRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// CheckExistence คือการจำลอง method CheckExistence
func (m *MockAccessControlRuleRepository) CheckExistence(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}
