package models

// UserLocation represents the location information of a user
type MessageBody struct {
	Body struct {
		Owner       string `json:"owner"`
		Reserver    string `json:"reserver"`
		Place       string `json:"place"`
		Date        string `json:"date"`
		NumOfPeople int    `json:"numOfPeople"`
		Payload     struct {
			TargetName string `json:"TargetName"`
			NewInfo    struct {
				AvailableSeat int `json:"AvailableSeat"`
			} `json:"NewInfo"`
		} `json:"payload"`
	} `json:"body"`
	Token string `json:"token"`
}
