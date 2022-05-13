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
	UserCreated = newEvent(users, user, create, false, 0)
	UserUpdated = newEvent(users, user, update, false, 0)
	UserDeleted = newEvent(users, user, delete, false, 0)

	verification          = "verification"
	UsersVerificationMail = newEvent(users, verification, mail, false, 0)
)

// Subscription events
var (
	subscription        = "subscription"
	SubscriptionCreated = newEvent(subscriptions, subscription, create, false, 0)
	SubscriptionUpdated = newEvent(subscriptions, subscription, update, false, 0)
	SubscriptionDeleted = newEvent(subscriptions, subscription, delete, false, 0)

	plan        = "plan"
	PlanCreated = newEvent(subscriptions, plan, create, false, 0)
	PlanUpdated = newEvent(subscriptions, plan, update, false, 0)
	PlanDeleted = newEvent(subscriptions, plan, delete, false, 0)
)

// Credential events
var (
	passwordReset                = "password_reset"
	CredentialsPasswordResetMail = newEvent(credentials, passwordReset, mail, false, 0)
)

// Core events
var (
	anchor          = "anchor"
	AnchorCreated   = newEvent(core, anchor, create, false, 0)
	AnchorUpdated   = newEvent(core, anchor, update, false, 0)
	AnchorFinalized = newEvent(core, anchor, "finalized", false, 0)

	anchorNetwork          = "anchor_network"
	AnchorNetworkCreated   = newEvent(core, anchorNetwork, create, false, 0)
	AnchorNetworkUpdated   = newEvent(core, anchorNetwork, update, false, 0)
	AnchorNetworkConfirmed = newEvent(core, anchorNetwork, "confirmed", false, 0)
)

// Transaction events
var (
	network          = "network"
	NetworkCreated   = newEvent(transactions, network, create, false, 0)
	NetworkConfirmed = newEvent(transactions, network, confirm, false, 0)

	anchorConfirm   = "anchor"
	AnchorConfirmed = newEvent(transactions, anchorConfirm, confirm, false, 0)

	anchorBloockchain           = "bloock_chain"
	AnchorBloockchainSender     = newEvent(transactions, anchorBloockchain, send, true, 30000)
	AnchorBloockchainSenderTest = newEvent(transactionsTest, anchorBloockchain, send, true, 30000)

	anchorRinkeby           = "ethereum_rinkeby"
	AnchorRinkebySender     = newEvent(transactions, anchorRinkeby, send, true, 30000)
	AnchorRinkebySenderTest = newEvent(transactionsTest, anchorRinkeby, send, true, 30000)

	anchorMainnet           = "ethereum_mainnet"
	AnchorMainnetSender     = newEvent(transactions, anchorMainnet, send, true, 30000)
	AnchorMainnetSenderTest = newEvent(transactionsTest, anchorMainnet, send, true, 30000)

	gnosisChain             = "gnosis_chain"
	AnchorGnosisChainSender = newEvent(transactions, gnosisChain, send, true, 30000)
)

// Notification events
var (
	webhook           = "webhook"
	WebhookScheduler  = newEvent(notifications, webhook, "schedule", false , 0)
	WebhookInvocation = newEvent(notifications, webhook, "invocation", false, 0)
	WebhookRetry      = newEvent(notifications, webhook, "retry", false, 0)
	WebhookConfirmed  = newEvent(notifications, webhook, "confirmed", false, 0)

	email     = "email"
	SendEmail = newEvent(notifications, email, "send", false, 0)
)

func newEvent(service, entity, action string, retry bool, expiration int) Type {
	return Type{
		name: fmt.Sprintf("%s.%s.%s", service, entity, action),
		retry: retry,
		expiration: expiration,
	}
}
