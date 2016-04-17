package uniq

import "github.com/sluu99/uuid"
import "crypto/md5"
import "encoding/hex"

func Uniq() string {

	uuid := uuid.Rand()

	h := md5.New()

	h.Write([]byte(uuid.Hex()))

	return hex.EncodeToString(h.Sum(nil))
}
