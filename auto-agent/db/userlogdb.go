package db

import (
	"encoding/csv"
	"io"
	"time"
)

// Generic message with no additional phase specific details
type UserLog struct {
	Username     string // Storing this in case if Users are deleted for some reasons
	Id           string
	PromptId     string
	PhaseId      string
	QuestionText string
	JsonResponse string
	ResponseId   string
	ResponseText string
	Mtype        string // ROBOT | HUMAN
	Date         time.Time
	URL          string
}

func WriteUserLogAsCSV(w io.Writer, log []UserLog) {
	csvWriter := csv.NewWriter(w)
	// Header
	csvWriter.Write([]string{
		"Username",
		"Id",
		"PromptId",
		"PhaseId",
		"QuestionText",
		"JsonResponse",
		"ResponseId",
		"ResponseText",
		"Mtype",
		"Date",
		"URL",
	})

	for _, row := range log {
		csvWriter.Write([]string{
			row.Username,
			row.Id,
			row.PromptId,
			row.PhaseId,
			row.QuestionText,
			row.JsonResponse,
			row.ResponseId,
			row.ResponseText,
			row.Mtype,
			row.Date.Format(time.RFC3339),
			row.URL,
		})
	}
	csvWriter.Flush()
}
