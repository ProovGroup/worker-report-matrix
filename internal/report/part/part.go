package part

import "worker-report-matrix/internal/permalink"

type Part struct {
	Title   string  `json:"title"`
	Note    string  `json:"note"`
	SignUrl string  `json:"sign_url,omitempty"`
	Infos   []*Info `json:"infos,omitempty"`
}

func (part *Part) ToPermalink() *Part {
	part.SignUrl = permalink.GetPermalink(part.SignUrl)
	return part
}
