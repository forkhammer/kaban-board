package models

type IssueType string
type IssueId string
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
	Assignees   []UserId
	WebUrl      string
	Labels      []LabelId
	ProjectId   ProjectId
	ReleaseId   *ReleaseId
	TaskType    *LabelId
	EstimateDev *uint
	EstimateQA  *uint
}
