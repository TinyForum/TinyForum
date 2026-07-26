package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateSlug() string {
	return fmt.Sprintf("%s-%s", time.Now().Format("20060102150405"), uuid.New().String()[:8])
}
