package auth

type Resource string

const (
	ResourceUsersMetadata             Resource = "users.metadata"
	ResourceUsersUser                 Resource = "users.user"
	ResourceUsersBusiness             Resource = "users.business"
	ResourceCredentialsApikey         Resource = "credentials.apikey"
	ResourceCredentialsTestApiKey     Resource = "credentials-test.apiKey"
	ResourceCredentialsSession        Resource = "credentials.session"
	ResourceCredentialsLicense        Resource = "credentials.license"
	ResourceCoreAnchor                Resource = "core.anchor"
	ResourceCoreProof                 Resource = "core.proof"
	ResourceCoreMessage               Resource = "core.message"
	ResourceCoreTestAnchor            Resource = "core-test.anchor"
	ResourceCoreTestProof             Resource = "core-test.proof"
	ResourceCoreTestMessage           Resource = "core-test.message"
	ResourceEventsActivity            Resource = "events.activity"
	ResourceEventsAnchor              Resource = "events.anchor"
	ResourceEventsWebhook             Resource = "events.webhook"
	ResourceEventsAnalytics           Resource = "events.analytics"
	ResourceSubscriptionsSubscription Resource = "subscriptions.subscription"
	ResourceSubscriptionsPlan         Resource = "subscriptions.plan"
	ResourceSubscriptionsInvoice      Resource = "subscriptions.invoice"
	ResourceNotificationsWebhook      Resource = "notifications.webhook"
	ResourceNotificationsFeedback     Resource = "notifications.feedback"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Ability struct {
	resource Resource
	action   Action
}

func (a Ability) Resource() string {
	return string(a.resource)
}

func (a Ability) Action() string {
	return string(a.action)
}

func NewAbility(resource Resource, action Action) Ability {
	return Ability{
		resource,
		action,
	}
}
