package uuid

import (
	"strings"

	"github.com/google/uuid"
)

type UID string

func NewUID(prefix string) UID {
	_uuid := uuid.New()
	__uid := strings.Join([]string{prefix, _uuid.String()}, "-")
	return UID(__uid)
}
