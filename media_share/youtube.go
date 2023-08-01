package media_share

import (
	"fmt"
	"github.com/youtube/media_share/dto"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"net/url"
	"path"
	"time"

	"net/http"
)

const (
	YoutubeFullURL  = "www.youtube.com"
	YoutubeShortURL = "www.youtu.be"
)

type YoutubeClient struct {
	Client *http.Client
}

func NewYoutubeClient(Key string) YoutubeClient {
	return YoutubeClient{
		Client: &http.Client{
			Transport: &transport.APIKey{
				Key: Key,
			},
		},
	}
}

func (c YoutubeClient) GetVideoInfo(videoURL string) (*dto.VideoInfoResponse, error) {
	service, err := youtube.New(c.Client)
	if err != nil {
		return nil, err
	}

	key, err := c.ParseVideoURL(videoURL)
	if err != nil {
		return nil, err
	}

	call := service.Videos.List([]string{
		"statistics",
		"snippet",
	}).Id(key)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) <= 0 {
		return nil, fmt.Errorf("statistic video not found")
	}

	item := resp.Items[0]
	published, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return nil, err
	}

	output := new(dto.VideoInfoResponse)
	output.VideoTitle = item.Snippet.Title
	output.ChannelTitle = item.Snippet.ChannelTitle
	output.ViewCount = item.Statistics.ViewCount
	output.CommentCount = item.Statistics.CommentCount
	output.DislikeCount = item.Statistics.DislikeCount
	output.FavoriteCount = item.Statistics.FavoriteCount
	output.LikeCount = item.Statistics.LikeCount
	output.PublishYear = published.Format("2006")

	return output, nil
}

func (c YoutubeClient) ParseVideoURL(videoURL string) (string, error) {
	u, err := url.Parse(videoURL)
	if err != nil {
		return "", err
	}

	var videoID string

	switch u.Host {
	case YoutubeFullURL:
		if path.Base(u.Path) == "watch" {
			m, _ := url.ParseQuery(u.RawQuery)
			key := m["v"]
			if key == nil {
				return "", fmt.Errorf("failed parsed url from youtube")
			}

			videoID = key[0]
		} else {
			videoID = path.Base(u.Path)
		}
	case YoutubeShortURL:
		videoID = path.Base(u.Path)
	}

	return videoID, nil
}
