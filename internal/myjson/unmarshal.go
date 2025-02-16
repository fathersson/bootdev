package myjson

import (
	"encoding/json"
	"log"
)

func Unmarshal(body []byte, dist any) {
	err := json.Unmarshal(body, &dist)
	if err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
	}
}
