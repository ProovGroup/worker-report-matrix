package process

import "worker-report-matrix/internal/permalink"

type Picture struct {
	Original  string `json:"original"`
	Thumbnail string `json:"thumbnail"`
}

func (picture *Picture) ToPermalink() *Picture {
	picture.Original = permalink.GetPermalink(picture.Original)
	picture.Thumbnail = permalink.GetPermalink(picture.Thumbnail)
	return picture
}
