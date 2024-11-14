package api

type HealthEnum string

const (
	Health_Normal              HealthEnum = "NORMAL_0"
	Health_Failure             HealthEnum = "FAILURE_1"
	Health_CheckFunction       HealthEnum = "CHECK_FUNCTION_2"
	Health_OffSpec             HealthEnum = "OFF_SPEC_3"
	Health_MaintenanceRequired HealthEnum = "MAINTENANCE_REQUIRED_4"
)

type Health struct {
	Health      HealthEnum `json:"Health"`
	HealthScore byte       `json:"HealthScore"`
}
