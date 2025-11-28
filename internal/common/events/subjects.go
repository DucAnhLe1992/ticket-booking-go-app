package events

// Subject defines event subjects used across services.
type Subject string

const (
	SubjectTicketCreated      Subject = "ticket:created"
	SubjectTicketUpdated      Subject = "ticket:updated"
	SubjectOrderCreated       Subject = "order:created"
	SubjectOrderCancelled     Subject = "order:cancelled"
	SubjectExpirationComplete Subject = "expiration:complete"
	SubjectPaymentCreated     Subject = "payment:created"
	SubjectUserCreated        Subject = "user:created"
)
