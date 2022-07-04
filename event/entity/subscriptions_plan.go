package event_entity

type SubscriptionsPlan struct {
	PlanID                 string `json:"plan_id"`
	PlanScope              string `json:"plan_scope"`
	License                string `json:"license"`
	MaxApiKeys             int    `json:"max_api_keys"`
	MaxSubscriptionRecords int    `json:"max_sub_records"`
	Analytics              bool   `json:"analytics"`
}

func NewPlanEventEntity(planID, planScope, license string, maxApiKeys, maxSubRecords int, analytics bool) SubscriptionsPlan {
	return SubscriptionsPlan{
		PlanID:                 planID,
		PlanScope:              planScope,
		License:                license,
		MaxApiKeys:             maxApiKeys,
		MaxSubscriptionRecords: maxSubRecords,
		Analytics:              analytics,
	}
}
