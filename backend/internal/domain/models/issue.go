package models

type IssueType string
type IssueId uint
type IssueIid string

const (
	IssueTypeEpic        IssueType = "EPIC"
	IssueTypeTask        IssueType = "TASK"
	IssueTypeIncident    IssueType = "INCIDENT"
	IssueTypeIssue       IssueType = "ISSUE"
	IssueTypeKeyResult   IssueType = "KEY_RESULT"
	IssueTypeObjective   IssueType = "OBJECTIVE"
	IssueTypeRequirement IssueType = "REQUIREMENT"
	IssueTypeTestCase    IssueType = "TEST_CASE"
	IssueTypeTicket      IssueType = "TICKET"
)

type Issue struct {
	Id          IssueId
	Iid         IssueIid
	Title       string
	IssueType   IssueType
	Assignees   []User
	WebUrl      string
	Labels      []Label
	Project     Project
	Release     *Release
	TaskType    *Label
	EstimateDev *uint
	EstimateQA  *uint
}

func (i *Issue) Validate() error {
	return nil
}
