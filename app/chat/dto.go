package chat

type createChatReq struct {
	SecondUserId string `json:"second_user_id"`
}

type createMessageReq struct {
	Text string `json:"text"`
}
