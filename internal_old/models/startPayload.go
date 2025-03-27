package models

type ActionType string

var GetGroupInfo ActionType = "getGroupInfo"
var GetKidInfo ActionType = "getKidInfo"

type StartPayload struct {
	Action  ActionType
	Payload []string
}
