package entity

import "encoding/json"

type Notify struct {
	Target  string
	Subject string
	Body    string
}

func (n *Notify) Encode() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Notify) Decode(p []byte) error {
	return json.Unmarshal(p, n)
}
