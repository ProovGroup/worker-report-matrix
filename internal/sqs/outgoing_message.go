package sqs

import (
	"encoding/json"
	"fmt"
	"os"
	"worker-report-matrix/internal/db"
	r "worker-report-matrix/internal/report"

	"github.com/ProovGroup/env"
)

var (
	PDF_SQS_REGION = os.Getenv("PDF_SQS_REGION")
	PDF_SQS_QUEUE  = os.Getenv("PDF_SQS_QUEUE")
)

type OutgoingMessage struct {
	Source      Source      `json:"Source"`
	Destination Destination `json:"Destination"`
	DataSource  DataSource  `json:"DataSource"`
}

func (message OutgoingMessage) Send(e *env.Env) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Set event region to match the queue region
	e.SQSRegion = PDF_SQS_REGION
	return e.SqsSend(PDF_SQS_QUEUE, string(jsonMessage))
}

func (message *OutgoingMessage) GetReport(e *env.Env, proovCode string) error {
	// Get report from database
	var report *r.Report
	report, err := db.GetReport(e, proovCode)
	if err != nil {
		return err
	}

	// Transform S3 links to permalink URLs
	report.ToPermalink()

	// Add geo location information
	if report.Geoloc == nil {
		report.Geoloc = &r.GeoLoc{}
	}
	report.Geoloc.Initialize(report.Latitude, report.Longitude)

	if report.Geoloc.Address == "" {
		fmt.Println("[INFO] Geoloc address not found, fetching from mapbox")
		report.Geoloc.GetAddressFromMapBox()
		if report.Geoloc.Address != "" {
			db.SaveGeoLoc(e, report.Geoloc)
			fmt.Println("[INFO] Geoloc saved to db")
		} else {
			fmt.Println("[WARN] Geoloc address not found")
		}
	}

	// Add report to message
	message.DataSource.Content = report
	message.Destination.Path.Key = fmt.Sprintf(message.Destination.Path.Key, report.ProovCode)

	return nil
}

type Source struct {
	Path    Path          `json:"Path"`
	Options SourceOptions `json:"Options"`
}

type Destination struct {
	Path    Path               `json:"Path"`
	Options DestinationOptions `json:"Options"`
}

type DataSource struct {
	Content *r.Report    `json:"Content,omitempty"`
	Options DataSourceOptions `json:"Options"`
}

type Path struct {
	Region string `json:"Region"`
	Bucket string `json:"Bucket"`
	Key    string `json:"Key"`
}

type SourceOptions struct {
	Type string `json:"Type"`
}

type DestinationOptions struct {
	PDFOptions interface{} `json:"PDFOptions"`
}

type DataSourceOptions struct {
	Type string `json:"Type"`
}
