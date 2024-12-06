package world

type NotificationEvent interface {
	isNotificationEvent()
}

type unimplementedNotificationEvent struct{}

func (*unimplementedNotificationEvent) isNotificationEvent() {}

type ChairNotificationEventMatching struct {
	ServerRideID string
	ChairNotificationEvent

	unimplementedNotificationEvent
}

type ChairNotificationEventCompleted struct {
	ServerRideID string
	ChairNotificationEvent

	unimplementedNotificationEvent
}

type ChairNotificationEvent struct {
	User        ChairNotificationEventUserPayload
	Pickup      Coordinate
	Destination Coordinate
}

type ChairNotificationEventUserPayload struct {
	ID   string
	Name string
}

type UserNotificationEventMatching struct {
	ServerRideID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventEnRoute struct {
	ServerRideID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventPickup struct {
	ServerRideID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventCarrying struct {
	ServerRideID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventArrived struct {
	ServerRideID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventCompleted struct {
	ServerRideID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEvent struct {
	Pickup      Coordinate
	Destination Coordinate
	Fare        int
	Chair       *UserNotificationEventChairPayload
}

type UserNotificationEventChairPayload struct {
	ID    string
	Name  string
	Model string
	Stats UserNotificationEventChairStatsPayload
}

type UserNotificationEventChairStatsPayload struct {
	TotalRidesCount    int
	TotalEvaluationAvg float64
}
