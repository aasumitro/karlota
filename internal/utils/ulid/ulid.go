package ulid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

func Generate(prefix string) string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	b := make([]byte, 16)
	binary.LittleEndian.PutUint64(b, uint64(t))
	_, err := rand.Read(b[6:])
	if err != nil {
		return ""
	}
	s := fmt.Sprintf("%x", b)
	return fmt.Sprintf("%s_%s%s%s%s%s",
		prefix, s[0:8], s[8:12], s[12:16], s[16:20], s[20:])
}
