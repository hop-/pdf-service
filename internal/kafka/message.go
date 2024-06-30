package kafka

import "encoding/json"

type Message struct {
	// Binary representation of json message
	message []byte
}

func NewMessage(msg any) (*Message, error) {
	message, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return &Message{message}, nil
}

// Get message as given interface
func (m *Message) Get(i any) error {
	return json.Unmarshal(m.message, i)
}
