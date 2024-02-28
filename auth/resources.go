package auth

type Resource string

const (
	ResourceUsersUser                 Resource = "users.user"
	ResourceCredentialsApikey         Resource = "credentials.api_key"
	ResourceCoreAnchor                Resource = "core.anchor"
	ResourceCoreProof                 Resource = "core.proof"
	ResourceCoreMessage               Resource = "core.message"
	ResourceCoreTestAnchor            Resource = "core-test.anchor"
	ResourceCoreTestProof             Resource = "core-test.proof"
	ResourceCoreTestMessage           Resource = "core-test.message"
	ResourceSubscriptionsSubscription Resource = "subscriptions.subscription"
	ResourceSubscriptionsUsage        Resource = "subscriptions.usage"
	ResourceSubscriptionsPlan         Resource = "subscriptions.plan"
	ResourceSubscriptionsInvoice      Resource = "subscriptions.invoice"
	ResourceNotificationsWebhook      Resource = "notifications.webhook"
	ResourceNotificationsFeedback     Resource = "notifications.feedback"
	ResourceDataAvailabilityHosting   Resource = "data-availability.hosting"
	ResourceDataAvailabilityIpfs      Resource = "data-availability.ipfs"
	ResourceNodeProxyRedirect         Resource = "node-proxy.redirect"
	ResourceKeysKey                   Resource = "keys.key"
	ResourceKeysCertificate           Resource = "keys.certificate"
	ResourceKeysSign                  Resource = "keys.sign"
	ResourceKeysVerify                Resource = "keys.verify"
	ResourceKeysEncrypt               Resource = "keys.encrypt"
	ResourceKeysDecrypt               Resource = "keys.decrypt"
	ResourceKeysAccessControl         Resource = "keys.access_control"
	ResourceIdentityIssuer            Resource = "identity.issuer"
	ResourceIdentitySchema            Resource = "identity.schema"
	ResourceIdentityCredential        Resource = "identity.credential"
	ResourceIdentityRevocation        Resource = "identity.revocation"
	ResourceCertifierProcess          Resource = "certifier.process"
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
