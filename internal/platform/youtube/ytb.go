package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type YouTubeService struct {
	apiKey string
}

type YouTubeVideoDetails struct {
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
	Duration     string `json:"duration"`
	PublishedAt  string `json:"published_at"`
}

func NewYouTubeService() *YouTubeService {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY is not set in environment variables")
	}

	return &YouTubeService{
		apiKey: apiKey,
	}
}

func (ytb *YouTubeService) GetVideoDetails(videoID string) (*YouTubeVideoDetails, error) {
	if ytb.apiKey == "" {
		return nil, errors.New("API key is missing")
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&part=snippet,contentDetails&key=%s", videoID, ytb.apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch video details, status code: %d", resp.StatusCode)
	}

	var data struct {
		Items []struct {
			Snippet struct {
				Title      string `json:"title"`
				Thumbnails struct {
					Default struct {
						URL string `json:"url"`
					} `json:"default"`
				} `json:"thumbnails"`
				PublishedAt string `json:"publishedAt"`
			} `json:"snippet"`
			ContentDetails struct {
				Duration string `json:"duration"`
			} `json:"contentDetails"`
		} `json:"items"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data.Items) == 0 {
		return nil, errors.New("video not found")
	}

	item := data.Items[0]

	duration := parseISO8601Duration(item.ContentDetails.Duration)

	publishedAt, err := formatPublishedAt(item.Snippet.PublishedAt)
	if err != nil {
		return nil, err
	}

	return &YouTubeVideoDetails{
		Title:        item.Snippet.Title,
		ThumbnailURL: item.Snippet.Thumbnails.Default.URL,
		Duration:     duration,
		PublishedAt:  publishedAt,
	}, nil
}

func parseISO8601Duration(isoDuration string) string {
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(isoDuration)

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// formatPublishedAt : Chuyển đổi ISO 8601 thành yyyy/MM/dd
func formatPublishedAt(isoDate string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		return "", err
	}
	return parsedTime.Format("2006/01/02"), nil
}
