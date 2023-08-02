package media_share

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestYoutubeClient_GetVideoInfo(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		videoUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Youtube get statistic",
			fields: fields{
				Key: "AIzaSyDFvszvjLdULpUw8pB7kDQPOJo6nfptccg",
			},
			args: args{
				videoUrl: "https://www.youtube.com/watch?v=K4DfmVY41x8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewYoutubeClient(tt.fields.Key)
			got, err := c.GetVideoInfo(tt.args.videoUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideoInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Got : %v", got)
		})
	}
}

func TestYoutubeClient_ParseVideoURL(t *testing.T) {
	type fields struct {
		Client *http.Client
	}
	type args struct {
		videoURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		loop    int
	}{
		{
			name:    "Youtube parse url watch",
			fields:  fields{Client: &http.Client{}},
			args:    args{videoURL: "https://www.youtube.com/watch?v=lWKsGfSf1U0"},
			want:    "lWKsGfSf1U0",
			wantErr: false,
			loop:    1,
		},
		{
			name:    "Youtube parse url embeed",
			fields:  fields{Client: &http.Client{}},
			args:    args{videoURL: "https://www.youtube.com/embeed/lWKsGfSf1U0"},
			want:    "lWKsGfSf1U0",
			wantErr: false,
			loop:    1,
		},
		{
			name:    "Youtube parse url embeed",
			fields:  fields{Client: &http.Client{}},
			args:    args{videoURL: "https://www.youtu.be/lWKsGfSf1U0"},
			want:    "lWKsGfSf1U0",
			wantErr: false,
			loop:    1,
		},
		{
			name:    "Youtube rate limited",
			fields:  fields{Client: &http.Client{}},
			args:    args{videoURL: "https://www.youtu.be/lWKsGfSf1U0"},
			want:    "lWKsGfSf1U0",
			wantErr: false,
			loop:    100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewYoutubeClient("AIzaSyDFvszvjLdULpUw8pB7kDQPOJo6nfptccg")

			result := struct {
				ErrorCount   int
				ErrorMessage string
			}{}

			f, err := os.Create(fmt.Sprintf("test-result-%s.txt", strings.ToLower(strings.ReplaceAll(tt.name, " ", "-"))))
			if err != nil {
				t.Logf("ERROR. failed create file. %v", err)
				return
			}

			defer f.Close()

			wg := &sync.WaitGroup{}

			for i := 1; i <= tt.loop; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					got, err := c.ParseVideoURL(tt.args.videoURL)
					if (err != nil) != tt.wantErr {
						t.Errorf("ParseVideoURL() error = %v, wantErr %v", err, tt.wantErr)
						result.ErrorCount = result.ErrorCount + 1
						result.ErrorMessage += err.Error()
						return
					}
					if got != tt.want {
						t.Errorf("ParseVideoURL() got = %v, want %v", got, tt.want)
						result.ErrorCount = result.ErrorCount + 1
						result.ErrorMessage += fmt.Sprintf("ParseVideoURL() got = %v, want %v", got, tt.want)
						return
					}

					resp, err := c.GetVideoInfo(tt.args.videoURL)
					if err != nil {
						result.ErrorCount = result.ErrorCount + 1
						result.ErrorMessage += err.Error()
						return
					}

					message := fmt.Sprintf("Video URL: %s\nExecution Time: %s\nVideo Title: %s\nChannel Title: %s\nView Count: %d\nPublished Year: %s\n\n",
						tt.args.videoURL, time.Now().Format("15:04:05.000"), resp.VideoTitle, resp.ChannelTitle, resp.ViewCount, resp.PublishYear)
					_, err = f.WriteString(message)
					if err != nil {
						t.Logf("ERROR. failed write file. %v", err)
						return
					}

					fmt.Printf("execution time: %s\n views count: %d\n", time.Now().Format("15:04:05.000"), resp.ViewCount)
				}()
			}
			wg.Wait()

			t.Logf("testing result.\n Error count: %d\nError messages: %s\n Testing count: %d",
				result.ErrorCount,
				result.ErrorMessage,
				tt.loop,
			)
		})
	}
}
