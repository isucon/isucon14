package world

type NotificationEvent interface {
	isNotificationEvent()
}

type unimplementedNotificationEvent struct{}

func (*unimplementedNotificationEvent) isNotificationEvent() {}

type ChairNotificationEventMatched struct {
	ServerRequestID string
	ChairNotificationEvent

	unimplementedNotificationEvent
}

type ChairNotificationEventDispatching struct {
	ServerRequestID string
	ChairNotificationEvent

	unimplementedNotificationEvent
}

type ChairNotificationEventDispatched struct {
	ServerRequestID string
	ChairNotificationEvent

	unimplementedNotificationEvent
}

type ChairNotificationEventCarrying struct {
	ServerRequestID string
	ChairNotificationEvent

	unimplementedNotificationEvent
}

type ChairNotificationEventArrived struct {
	ServerRequestID string
	ChairNotificationEvent

	unimplementedNotificationEvent
}

type ChairNotificationEventCompleted struct {
	ServerRequestID string
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
	ServerRequestID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventDispatching struct {
	ServerRequestID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventDispatched struct {
	ServerRequestID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventCarrying struct {
	ServerRequestID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventArrived struct {
	ServerRequestID string
	UserNotificationEvent

	unimplementedNotificationEvent
}

type UserNotificationEventCompleted struct {
	ServerRequestID string
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
