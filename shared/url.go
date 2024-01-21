package shared

import (
	"fmt"
	"os"
)

func GetShortenedUrl(slug string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), slug)
}
