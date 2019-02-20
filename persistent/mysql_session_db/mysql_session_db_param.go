package mysql_session_db

type CAddSessionInput struct {
	SessionUuid string
	SessionUuidIsValid bool
	TimeoutS int64
	TimeoutSIsValid bool
	LoseValidTime int64
	LoseValidTimeIsValid bool
}

type CDeleteSessionInput struct {
	SessionUuid string
	SessionUuidIsValid bool
}

type CUpdateSessionInput struct {
	Condition string
	ConditionIsValid bool
	SessionUuid string
	SessionUuidIsValid bool
}

type CGetSessionInput struct {
	SessionUuid string
	SessionUuidIsValid bool
}

type CGetSessionOutput struct {
	TimeoutS int64
	TimeoutSIsValid bool
	LoseValidTime int64
	LoseValidTimeIsValid bool
}

type CGetCountBySessionUuidInput struct {
	SessionUuid string
	SessionUuidIsValid bool
}

type CGetCountBySessionUuidOutput struct {
	Count int
	CountIsValid bool
}

