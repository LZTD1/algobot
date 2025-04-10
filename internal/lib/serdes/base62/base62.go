package base62

import (
	"algobot/internal/domain"
	"fmt"
	"github.com/jxskiss/base62"
	"log/slog"
	"strconv"
	"strings"
)

type Serdes struct {
	log *slog.Logger
}

func NewSerdes(log *slog.Logger) *Serdes {
	return &Serdes{log: log}
}

func (s *Serdes) Serialize(msg domain.SerializeMessage) (string, error) {
	encoded := base62.EncodeToString([]byte(fmt.Sprintf(
		"%d-%s",
		msg.Type,
		strings.Join(msg.Data, ","),
	)))
	return encoded, nil
}

func (s *Serdes) Deserialize(decoded string) (*domain.SerializeMessage, error) {
	const op = "serdes.Deserialize"
	log := s.log.With(
		slog.String("op", op),
	)

	encoded, err := base62.DecodeString(decoded)
	if err != nil {
		log.Warn("Failed to decode serdes")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	encodedMsg := strings.Split(string(encoded), "-")
	encodedType := encodedMsg[0]
	encodedData := strings.Split(encodedMsg[1], ",")

	encodedID, err := strconv.Atoi(encodedType)
	if err != nil {
		log.Warn("Failed to decode serdes")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	serType := domain.SerType(encodedID)

	return &domain.SerializeMessage{
		Type: serType,
		Data: encodedData,
	}, nil
}
