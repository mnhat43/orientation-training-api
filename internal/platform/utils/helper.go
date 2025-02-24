package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"
)

func GetSHA256Hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetKeyToken() string {
	keyToken := os.Getenv("KEY_TOKEN")

	return keyToken
}

func TimeNowUTC() time.Time {
	utc, _ := time.LoadLocation("UTC")
	return time.Now().In(utc)
}

func FindIntInSlice(slice []int, val int) bool {
	if len(slice) == 0 {
		return false
	}

	for _, item := range slice {
		if item == val {
			return true
		}
	}

	return false
}

// FindStringInArray : find item in array
// Params    : array, string item
// Returns   : index, bool
func FindStringInArray(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// func parseISO8601Duration(isoDuration string) string {
// 	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
// 	matches := re.FindStringSubmatch(isoDuration)

// 	hours, _ := strconv.Atoi(matches[1])
// 	minutes, _ := strconv.Atoi(matches[2])
// 	seconds, _ := strconv.Atoi(matches[3])

// 	if hours > 0 {
// 		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
// 	}
// 	return fmt.Sprintf("%02d:%02d", minutes, seconds)
// }

// // formatPublishedAt : Chuyển đổi ISO 8601 thành yyyy/MM/dd
// func formatPublishedAt(isoDate string) (string, error) {
// 	parsedTime, err := time.Parse(time.RFC3339, isoDate)
// 	if err != nil {
// 		return "", err
// 	}
// 	return parsedTime.Format("2006/01/02"), nil
// }
