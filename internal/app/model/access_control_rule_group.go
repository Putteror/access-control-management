package model

type AccessControlRuleGroup struct {
	BaseModel
	AccessControlGroupID string `json:"access_control_group_id"`
	AccessControlRuleID  string `json:"access_control_rule_id"`
}
