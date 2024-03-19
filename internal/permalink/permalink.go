package permalink

import (
	"fmt"
	"os"

	permalink "github.com/ProovGroup/lib-permalink"
)

var (
	REGION = os.Getenv("PICTURES_REGION")
	BUCKET = os.Getenv("PICTURES_BUCKET")
	ENV	= os.Getenv("ENV")
)

func GetPermalink(key string) string {
	if (key == "") {
		return ""
	}

	baseURL := "https://permalink.weproov.com"
	link := permalink.NewPermalink(ENV, permalink.S3, 0)

	url, err := link.SetRegion(REGION).SetBucket(BUCKET).SetKey(key).AppendToURL(baseURL)
	if err != nil {
		fmt.Println("[ERROR] GetPermalink:", err)
		return ""
	}

	return url
}


