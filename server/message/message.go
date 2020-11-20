package message

type Message struct {
	User     string
	Message  string
	Type     int
	File     []byte
	FileName string
}
