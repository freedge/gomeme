// Package types define various structures used to interact with Control-M API
package types

import "time"

// QR contains the definition of a quantitative resource
type QR struct {
	Name      string
	Ctm       string
	Available string // Available is an int but defined as a String
	Max       int
}

// SessionLoginQuery contains the data to pass when login in
type SessionLoginQuery struct {
	Username string
	Password string
}

// Token is contained in the login response
type Token struct {
	Version  string
	Token    string
	Username string
}

// Status is a job status
type Status struct {
	JobId          string
	FolderID       string `json:"folderId"`
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

// JobsStatusReply is the reply from run/jobs/status
type JobsStatusReply struct {
	Statuses []Status
	Returned int
	Total    int
}

// ParseTime parses a non empty StartTime or EndTime from a Status structure and returns a Time
func ParseTime(s string) (t time.Time, err error) {
	return time.Parse("20060102150405", s)
}

// OrderQuery contains the data to order a new job
type OrderQuery struct {
	Hold   bool
	Ctm    string
	Folder string
	Jobs   string
}

// SetResourceQuery contain the data to set a QR. QR name is given in the URL parameter
type SetResourceQuery struct {
	Max string
}

// Message is returned by various services
type Message struct {
	Message string
}

// SetResourceReply is returned when setting a QR
type SetResourceReply = Message

// JobActionReply is returned for most actions
type JobActionReply = Message

// OrderJobReply is returned after ordering a job
type OrderJobReply struct {
	RunID     string `json:"runId"` // this syntax should probably be adopted everywhere in the file
	StatusURI string
}

// ErrorReply is returned in case of errors most of the time
type ErrorReply struct {
	Errors []Message
}

// Server replied by the config.servers service
type Server struct {
	Name    string
	Host    string
	State   string
	Message string
}

// ConfigServersReply is the reply to the config server
type ConfigServersReply []Server

// Agent replied by the config.agents service
type Agent struct {
	NodeID string
	Status string
}

// ConfigAgentsReply is the reply to the config.agents service
type ConfigAgentsReply struct {
	Agents []Agent
}

// ConfigAgentParam is one paramater of the agent
type ConfigAgentParam struct {
	Name, Value, DefaultValue string
}

// ConfigAgentParamsReply is the reply to the config agent params action
type ConfigAgentParamsReply []ConfigAgentParam

// LogoutReply is the message returned after login out
type LogoutReply = Message
