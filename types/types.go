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
	JobId              string
	FolderID           string `json:"folderId"`
	NumberOfRuns       int
	Name               string
	Folder             string
	Type               string
	Status             string
	Held               bool
	Deleted            bool
	StartTime          string
	EndTime            string
	OrderDate          string
	Ctm                string
	Description        string
	Host               string
	Application        string
	SubApplication     string
	OutputURI          string
	LogURI             string
	EstimatedStartTime []string
	EstimatedEndTime   []string
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

// ParseTime2 parses a non empty StartTime or EndTime from a Status structure and returns a Time.
// Actually needs a second one since the time is different when retrieving through job.status
func ParseTime2(s string) (t time.Time, err error) {
	return time.Parse("Jan 2, 2006 3:04:05 PM", s)
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

// AddResourceQuery contain the data to create a QR
type AddResourceQuery struct {
	Max  string `json:"max"`
	Name string `json:"name"`
}

// Message is returned by various services
type Message struct {
	Message string
	File    string
	Line    int
	Col     int
	ID      string `json:"id"`
	Item    string
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
	Version string
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

// PingAgentReply is the message returned when pinging an agent
type PingAgentReply = Message

// PingAgentQuery are the parameters sent to ping an agent
type PingAgentQuery struct {
	Timeout  int  `json:"timeout"`
	Discover bool `json:"discover"`
}

// DeployErrors are the kind of errors you get when deploying
type DeployErrors struct {
	Lines []string
}

// DeployResult is the result of a deployment
type DeployResult struct {
	DeploymentFile                    string
	SuccessfulFoldersCount            int
	SuccessfulSmartFoldersCount       int
	SuccessfulSubFoldersCount         int
	SuccessfulJobsCount               int
	SuccessfulConnectionProfilesCount int
	SuccessfulDriversCount            int
	IsDeployDescriptorValid           bool
	DeployedFolders                   []string
	Errors                            []DeployErrors
}

// DeployReply is the reply to the deployment service
type DeployReply = []DeployResult

// HostInGroup is one host part of a host group
type HostInGroup struct {
	Host string
}

// HostGroupAgentsReply is the reply to the hostgroup/agents service
type HostGroupAgentsReply = []HostInGroup

// OneJobDefinition is part of one job replied by JobGetReply
type OneJobDefinition struct {
	FilePath       string
	FileName       string
	Description    string
	RunAs          string
	Application    string
	SubApplication string
	Arguments      []string
}

// JobGetReply is what the job/get service will reply
type JobGetReply = map[string]OneJobDefinition

// WaitinfoReply is the information from the waiting info
type WaitinfoReply = []string

// OneFolderDefinition is a folder containing a few jobs
type OneFolderDefinition struct {
	ControlmServer string
	SiteStandard   string
	Type           string
	OrderMethod    string
}

// DeployPutFormat is part of what you give in input to the deploy or build service
type DeployPutFormat = map[string]OneFolderDefinition
