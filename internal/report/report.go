package report

import (
	part "worker-report-matrix/internal/report/part"
	process "worker-report-matrix/internal/report/process"
)

type Report struct {
	ProovCode      string             `json:"proov_code"`
	Owner          int                `json:"owner"`
	Title          string             `json:"title"`
	IdentifierItem string             `json:"identifier_item"`
	Logo           string             `json:"logo"`
	CreatedAt      string             `json:"created_at"`
	Latitude       float64            `json:"latitude"`
	Longitude      float64            `json:"longitude"`
	Geoloc         *GeoLoc            `json:"geoloc,omitempty"`
	Parts          []*part.Part       `json:"parts,omitempty"`
	Processes      []*process.Process `json:"processes,omitempty"`
}

func (report *Report) ToPermalink() *Report {
	for i := 0; i < len(report.Processes); i++ {
		report.Processes[i].ToPermalink()
	}
	for i := 0; i < len(report.Parts); i++ {
		report.Parts[i].ToPermalink()
	}
	return report
}
