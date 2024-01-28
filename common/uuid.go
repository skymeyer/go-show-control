package common

import (
	"github.com/google/uuid"
)

func StableUUIDString(in string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(in)).String()
}
