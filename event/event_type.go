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
	mail   = "mail"
)

// User events
var (
	user        = "user"
	UserCreated = newEvent(users, user, create)
	UserUpdated = newEvent(users, user, update)
	UserDeleted = newEvent(users, user, delete)

	verification          = "verification"
	UsersVerificationMail = newEvent(users, verification, mail)
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

// Credential events
var (
	passwordReset                = "password_reset"
	CredentialsPasswordResetMail = newEvent(credentials, passwordReset, mail)
)

// Core events
var (
	anchor          = "anchor"
	AnchorCreated   = newEvent(core, anchor, create)
	AnchorUpdated   = newEvent(core, anchor, update)
	AnchorFinalized = newEvent(core, anchor, "finalized")

	anchorNetwork          = "anchor_network"
	AnchorNetworkCreated   = newEvent(core, anchorNetwork, create)
	AnchorNetworkUpdated   = newEvent(core, anchorNetwork, update)
	AnchorNetworkConfirmed = newEvent(core, anchorNetwork, "confirmed")
)

// Transaction events
var (
	network          = "network"
	NetworkCreated   = newEvent(transactions, network, create)
	NetworkConfirmed = newEvent(transactions, network, "confirmed")
)

// Notification events
var (
	webhook           = "webhook"
	WebhookScheduler  = newEvent(notifications, webhook, "schedule")
	WebhookInvocation = newEvent(notifications, webhook, "invocation")
	WebhookRetry      = newEvent(notifications, webhook, "retry")
	WebhookConfirmed  = newEvent(notifications, webhook, "confirmed")

	email     = "email"
	SendEmail = newEvent(notifications, email, "send")
)

func newEvent(service, entity, action string) Type {
	return Type(fmt.Sprintf("%s.%s.%s", service, entity, action))
}
