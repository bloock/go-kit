package domain

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
	node             = "node"
	dataAvailability = "data-availability"
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

// Events events
var (
	activity                   = "activity"
	segment                    = "segment"
	EventsActivityCreated      = newEvent(events, activity, create, newEventArgs{})
	EventsSegmentEventReceived = newEvent(events, segment, create, newEventArgs{})
)

// Subscription events
var (
	subscription        = "subscription"
	SubscriptionCreated = newEvent(subscriptions, subscription, create, newEventArgs{})
	SubscriptionUpdated = newEvent(subscriptions, subscription, update, newEventArgs{})
	SubscriptionDeleted = newEvent(subscriptions, subscription, delete, newEventArgs{})

	product        = "product"
	ProductCreated = newEvent(subscriptions, product, create, newEventArgs{})
	ProductUpdated = newEvent(subscriptions, product, update, newEventArgs{})
	ProductDeleted = newEvent(subscriptions, product, delete, newEventArgs{})
)

// Credential events
var (
	passwordReset                = "password_reset"
	CredentialsPasswordResetMail = newEvent(credentials, passwordReset, mail, newEventArgs{})
)

// Core events
var (
	anchorNetwork          = "anchor_network"
	AnchorNetworkConfirmed = newEvent(core, anchorNetwork, "confirmed", newEventArgs{})

	recordCounter        = "record_counter"
	RecordCounterCreated = newEvent(core, recordCounter, create, newEventArgs{})

	usageRecordLimit        = "usage_record_limit"
	UsageRecordLimitUpdated = newEvent(core, usageRecordLimit, update, newEventArgs{})

	usageRecord        = "usage_record"
	UsageRecordUpdated = newEvent(core, usageRecord, update, newEventArgs{})
	UsageRecordDeleted = newEvent(core, usageRecord, delete, newEventArgs{})
)

// Node events
var (
	requestCounter        = "request_counter"
	RequestCounterCreated = newEvent(node, requestCounter, create, newEventArgs{})

	usageRequestLimit        = "usage_request_limit"
	UsageRequestLimitUpdated = newEvent(node, usageRequestLimit, update, newEventArgs{})

	usageRequest        = "usage_request"
	UsageRequestUpdated = newEvent(node, usageRequest, update, newEventArgs{})
	UsageRequestDeleted = newEvent(node, usageRequest, delete, newEventArgs{})
)

// Data Availability events
var (
	storageWorker        = "storage_worker"
	StorageWorkerCreated = newEvent(dataAvailability, storageWorker, create, newEventArgs{})

	transferWorker        = "transfer_worker"
	TransferWorkerCreated = newEvent(dataAvailability, transferWorker, create, newEventArgs{})

	usageStorageLimit        = "usage_storage_limit"
	UsageStorageLimitUpdated = newEvent(dataAvailability, usageStorageLimit, update, newEventArgs{})

	usageTransferLimit        = "usage_transfer_limit"
	UsageTransferLimitUpdated = newEvent(dataAvailability, usageTransferLimit, update, newEventArgs{})

	usageStorage        = "usage_storage"
	UsageStorageUpdated = newEvent(dataAvailability, usageStorage, update, newEventArgs{})
	UsageStorageDeleted = newEvent(dataAvailability, usageStorage, delete, newEventArgs{})

	usageTransfer        = "usage_transfer"
	UsageTransferUpdated = newEvent(dataAvailability, usageTransfer, update, newEventArgs{})
	UsageTransferDeleted = newEvent(dataAvailability, usageTransfer, delete, newEventArgs{})
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

	email     = "email"
	SendEmail = newEvent(notifications, email, "send", newEventArgs{})
)

type newEventArgs struct {
	retry      bool
	expiration int
}

func newEvent(service, entity, action string, args newEventArgs) EventType {
	return EventType{
		name:       fmt.Sprintf("%s.%s.%s", service, entity, action),
		retry:      args.retry,
		expiration: args.expiration,
	}
}
