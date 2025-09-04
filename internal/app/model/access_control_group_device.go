package model

type AccessControlGroupDevice struct {
	BaseModel
	AccessControlGroupID  string `json:"access_control_group_id"`
	AccessControlDeviceID string `json:"access_control_device_id"`
}
