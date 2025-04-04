package base62

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/serdes"
	"fmt"
	"github.com/jxskiss/base62"
	"log/slog"
)

type Serdes struct {
	log *slog.Logger
}

func NewSerdes(log *slog.Logger) *Serdes {
	return &Serdes{log: log}
}

func (s *Serdes) Serialize(group models.Group, traceID interface{}) (string, error) {
	encoded := base62.EncodeToString([]byte(fmt.Sprintf(
		"%d-%d",
		serdes.GroupType,
		group.GroupID,
	)))
	return encoded, nil
}
