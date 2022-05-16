package viewmodels

type UDPMessage struct {
	ID      int    `json:"id"`
	Message string `json:"msg"`
	User    string `json:"user"`
}

type UDPBroadcastMessage struct {
	DestinationAddress string `json:"-"`
	Message            string
}
