package dto

type VideoInfoResponse struct {
	CommentCount  uint64 `json:"commentCount,omitempty,string"`
	DislikeCount  uint64 `json:"dislikeCount,omitempty,string"`
	FavoriteCount uint64 `json:"favoriteCount,omitempty,string"`
	LikeCount     uint64 `json:"likeCount,omitempty,string"`
	ViewCount     uint64 `json:"viewCount,omitempty,string"`
	VideoTitle    string `json:"videoTitle,omitempty,string"`
	ChannelTitle  string `json:"channelTitle,omitempty,string"`
	PublishYear   string `json:"publishYear,omitempty,string"`
}
