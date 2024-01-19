package utils

import "github.com/google/uuid"

func ParserUUID(value string) (id uuid.UUID, err error) {
	id, err = uuid.Parse(value)
	if err != nil {
		return id, err
	}
	return id, nil
}
