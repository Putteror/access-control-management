package schema

type AccessControlRuleSearchQuery struct {
	Name  string `form:"name"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	All   bool   `form:"all"`
}

type AccessControlRuleRequest struct {
	Name                  string   `form:"name" validate:"required"`
	AccessControlGroupIDs []string `json:"accessControlGroupIds"`
}

type AccessControlRuleInfoResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AccessControlRuleResponse struct {
	ID                  string                           `json:"id"`
	Name                string                           `json:"name"`
	AccessControlGroups []AccessControlGroupInfoResponse `json:"accessControlGroups"`
}
