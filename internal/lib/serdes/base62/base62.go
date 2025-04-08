package base62

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/serdes"
	"fmt"
	"github.com/jxskiss/base62"
	"log/slog"
	"strings"
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

func (s *Serdes) GetType(decoded string) (serdes.SerType, error) {
	const op = "serdes.GetType"
	log := s.log.With(
		slog.String("op", op),
	)

	encoded, err := base62.DecodeString(decoded)
	if err != nil {
		log.Warn("Failed to decode serdes")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	encodedType := strings.Split(string(encoded), "-")[0]

	switch encodedType {
	case "0":
		return serdes.GroupType, nil
	case "1":
		return serdes.UserType, nil
	default:
		return 0, fmt.Errorf("%s is not a recognized serdes : %w", op, serdes.ErrUnrecognized)
	}
}
