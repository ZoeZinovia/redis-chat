package viewmodels

type UDPMessage struct {
	ID      int    `json:"id"`
	Message string `json:"msg"`
	User    string `json:"user"`
}

type UDPMessages struct {
	Messages []UDPMessage `json:"messages"`
}
