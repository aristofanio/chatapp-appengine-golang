package uuid

import (
	"strings"

	"github.com/google/uuid"
)

type UID string

func (u UID) Equals(oth UID) bool {
	return string(u) == string(oth)
}

func NewUID(prefix string) UID {
	_uuid := uuid.New()
	__uid := strings.Join([]string{prefix, _uuid.String()}, "-")
	return UID(__uid)
}
