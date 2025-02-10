package serdes

import (
	"fmt"
	"github.com/jxskiss/base62"
	"strings"
	"tgbot/internal/models"
)

func Serialize(m models.StartPayload) string {
	return base62.EncodeToString([]byte(fmt.Sprintf("%s-%s", m.Action, strings.Join(m.Payload, "-"))))
}
func Deserialize(s string) (models.StartPayload, error) {
	decodeString, err := base62.DecodeString(s)
	if err != nil {
		return models.StartPayload{}, err
	}
	arr := strings.Split(string(decodeString), "-")

	m := models.StartPayload{
		Action:  models.ActionType(arr[0]),
		Payload: arr[1:],
	}
	return m, nil
}
