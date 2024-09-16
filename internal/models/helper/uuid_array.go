package helper

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (u *UUIDArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("не удалось сканировать поле approved_users: неверный тип данных")
	}

	return json.Unmarshal(bytes, u)
}

func (u UUIDArray) Value() (driver.Value, error) {
	return json.Marshal(u)
}
