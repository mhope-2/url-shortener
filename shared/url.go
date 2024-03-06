// Package shared houses shared components
package shared

import (
	"fmt"
	"os"
)

func GetShortenedURL(slug string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), slug)
}
