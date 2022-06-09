package event

import "fmt"

// Services
var (
	users            = "users"
	subscriptions    = "subscriptions"
	credentials      = "credentials"
	events           = "events"
	notifications    = "notifications"
	transactions     = "transactions"
	transactionsTest = "transactions-test"
	core             = "core"
)

// Actions
var (
	create  = "create"
	update  = "update"
	delete  = "delete"
	mail    = "mail"
	send    = "send"
	confirm = "confirmed"
)

// User events
var (
	user        = "user"
	UserCreated = newEvent(users, user, create, newEventArgs{})
	UserUpdated = newEvent(users, user, update, newEventArgs{})
	UserDeleted = newEvent(users, user, delete, newEventArgs{})

	verification          = "verification"
	UsersVerificationMail = newEvent(users, verification, mail, newEventArgs{})
)

// Subscription events
var (
	subscription        = "subscription"
	SubscriptionCreated = newEvent(subscriptions, subscription, create, newEventArgs{})
	SubscriptionUpdated = newEvent(subscriptions, subscription, update, newEventArgs{})
	SubscriptionDeleted = newEvent(subscriptions, subscription, delete, newEventArgs{})

	plan        = "plan"
	PlanCreated = newEvent(subscriptions, plan, create, newEventArgs{})
	PlanUpdated = newEvent(subscriptions, plan, update, newEventArgs{})
	PlanDeleted = newEvent(subscriptions, plan, delete, newEventArgs{})
)

// Credential events
var (
	passwordReset                = "password_reset"
	CredentialsPasswordResetMail = newEvent(credentials, passwordReset, mail, newEventArgs{})
)

// Core events
var (
	anchor          = "anchor"
	AnchorCreated   = newEvent(core, anchor, create, newEventArgs{})
	AnchorUpdated   = newEvent(core, anchor, update, newEventArgs{})
	AnchorFinalized = newEvent(core, anchor, "finalized", newEventArgs{})

	anchorNetwork          = "anchor_network"
	AnchorNetworkCreated   = newEvent(core, anchorNetwork, create, newEventArgs{})
	AnchorNetworkUpdated   = newEvent(core, anchorNetwork, update, newEventArgs{})
	AnchorNetworkConfirmed = newEvent(core, anchorNetwork, "confirmed", newEventArgs{})
)

// Transaction events
var (
	network          = "network"
	NetworkCreated   = newEvent(transactions, network, create, newEventArgs{})
	NetworkConfirmed = newEvent(transactions, network, confirm, newEventArgs{})

	anchorConfirm   = "anchor"
	AnchorConfirmed = newEvent(transactions, anchorConfirm, confirm, newEventArgs{})
)

// Notification events
var (
	webhook           = "webhook"
	WebhookScheduler  = newEvent(notifications, webhook, "schedule", newEventArgs{})
	WebhookInvocation = newEvent(notifications, webhook, "invocation", newEventArgs{})
	WebhookRetry      = newEvent(notifications, webhook, "retry", newEventArgs{})
	WebhookConfirmed  = newEvent(notifications, webhook, "confirmed", newEventArgs{})

	email     = "email"
	SendEmail = newEvent(notifications, email, "send", newEventArgs{})
)

type newEventArgs struct {
	retry      bool
	expiration int
}

func newEvent(service, entity, action string, args newEventArgs) Type {
	return Type{
		name:       fmt.Sprintf("%s.%s.%s", service, entity, action),
		retry:      args.retry,
		expiration: args.expiration,
	}
}
