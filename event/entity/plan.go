package event_entity

type Plan struct {
	PlanID                 string `json:"plan_id"`
	PlanScope              string `json:"plan_scope"`
	License                string `json:"license"`
	MaxApiKeys             int    `json:"max_api_keys"`
	MaxSubscriptionRecords int    `json:"max_sub_records"`
}

func NewPlanEventEntity(planID, planScope, license string, maxApiKeys, maxSubRecords int) Plan {
	return Plan{
		PlanID:                 planID,
		PlanScope:              planScope,
		License:                license,
		MaxApiKeys:             maxApiKeys,
		MaxSubscriptionRecords: maxSubRecords,
	}
}
