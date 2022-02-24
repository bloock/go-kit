package event

import "fmt"

// Services
var (
	users         = "users"
	subscriptions = "subscriptions"
	credentials   = "credentials"
	events        = "events"
	notifications = "notifications"
	transactions  = "transactions"
	core          = "core"
)

// Actions
var (
	create = "create"
	update = "update"
	delete = "delete"
)

// User events
var (
	user        = "user"
	UserCreated = newEvent(users, user, create)
	UserUpdated = newEvent(users, user, update)
	UserDeleted = newEvent(users, user, delete)
)

// Subscription events
var (
	subscription        = "subscription"
	SubscriptionCreated = newEvent(subscriptions, subscription, create)
	SubscriptionUpdated = newEvent(subscriptions, subscription, update)
	SubscriptionDeleted = newEvent(subscriptions, subscription, delete)

	plan        = "plan"
	PlanCreated = newEvent(subscriptions, plan, create)
	PlanUpdated = newEvent(subscriptions, plan, update)
	PlanDeleted = newEvent(subscriptions, plan, delete)
)

// Core events
var (
	anchor          = "anchor"
	AnchorCreated   = newEvent(core, anchor, create)
	AnchorFinalized = newEvent(core, anchor, "finalized")
)

// Transaction events
var (
	network          = "network"
	NetworkCreated   = newEvent(transactions, network, create)
	NetworkConfirmed = newEvent(transactions, network, "confirmed")
)

func newEvent(service, entity, action string) Type {
	return Type(fmt.Sprintf("%s.%s.%s", service, entity, action))
}
