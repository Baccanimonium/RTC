package RTC

import "encoding/json"

const (
	BroadcastDeleteChatMessage = "onDeleteChatMessage"
	BroadcastUpdateChatMessage = "onUpdateChatMessage"
	BroadcastCreateChatMessage = "onNewChatMessage"
	BroadcastCreateChannel     = "onNewChannel"
	BroadcastDeleteChannel     = "onDeleteChannel"
	BroadcastCreateSchedule    = "onNewSchedule"
	BroadcastUpdateSchedule    = "onUpdateSchedule"
	BroadcastDeleteSchedule    = "onDeleteSchedule"
	RTCCandidate               = "candidate"
	RTCAnswer                  = "answer"
)

func ConvertToJson(value interface{}) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(value)
	err := json.Unmarshal(inrec, &inInterface)

	return inInterface, err
}
