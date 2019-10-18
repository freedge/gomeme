// Package types define various structures used to interact with Control-M API
package types

import "time"

type QR struct {
	Name      string
	Ctm       string
	Available string // Available is an int but defined as a String
	Max       int
}

type SessionLoginQuery struct {
	Username string
	Password string
}

type Token struct {
	Version  string
	Token    string
	Username string
}

type Status struct {
	JobId          string
	FolderId       string
	NumberOfRuns   int
	Name           string
	Folder         string
	Type           string
	Status         string
	Held           bool
	Deleted        bool
	StartTime      string
	EndTime        string
	OrderDate      string
	Ctm            string
	Description    string
	Host           string
	Application    string
	SubApplication string
	OutputURI      string
	LogURI         string
}

type JobsStatusReply struct {
	Statuses []Status
	Returned int
	Total    int
}

// ParseTime parses a non empty StartTime or EndTime from a Status structure and returns a Time
func ParseTime(s string) (t time.Time, err error) {
	t, err = time.Parse("20060102150405", s)
	return
}

type OrderQuery struct {
	Hold   bool
	Ctm    string
	Folder string
	Jobs   string
}
