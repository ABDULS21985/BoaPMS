package email

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"
)

// TemplateData holds all possible placeholder values for email templates.
// Each field maps to a .NET HtmlMessageTemplate placeholder:
//
//	{recipient}      → Recipient
//	{requester}      → Requester
//	{requestor}      → Requester  (alias used in some .NET templates)
//	{requestorName}  → RequestorName
//	{app_url}        → AppURL
//	{CompetencyInfo} → CompetencyInfo
//	{reviewperiod}   → ReviewPeriod
//	{deadlinedate}   → DeadlineDate
//	{competencyname} → CompetencyName
//	{description}    → Description
type TemplateData struct {
	Recipient      string
	Requester      string
	RequestorName  string
	AppURL         string
	CompetencyInfo string
	ReviewPeriod   string
	DeadlineDate   string
	CompetencyName string
	Description    string
}

// Template keys used in Render(). Each corresponds to a .NET HtmlMessageTemplate field.
const (
	TplNewCompetencyToApprover       = "NewCompetencyToApprover"
	TplNewCompetencyToRequestor      = "NewCompetencyToRequestor"
	TplApprovedCompetencyToRequestor = "ApprovedCompetencyToRequestor"
	TplApprovedCompetencyToApprover  = "ApprovedCompetencyToApprover"
	TplNewReviewPeriodToRequestor    = "NewReviewPeriodToRequestor"
	TplNewReviewPeriodToApprover     = "NewReviewPeriodToApprover"
	TplApprovedReviewPeriodToAll     = "ApprovedReviewPeriodToAllStaff"
	TplReminderMessage               = "ReminderMessage"
	TplSelfReviewCompleted           = "SelfReviewCompleted"
	TplPeersReviewCompleted          = "PeersReviewCompleted"
	TplSubordinateReviewCompleted    = "SubordinateReviewCompleted"
	TplSuperiorReviewCompleted       = "SuperiorReviewCompleted"
	TplDevTaskAssigned               = "DevTaskAssigned"
	TplDevTaskCompleted              = "DevTaskCompleted"
	TplDevTaskApproved               = "DevTaskApproved"
	TplJobRoleUpdateRequest          = "JobRoleUpdateRequest"
	TplApprovedJobRole               = "ApprovedJobRole"
	TplRejectJobRole                 = "RejectJobRole"
	TplCMSTeamApprovedJobRole        = "CMSTeamApprovedJobRole"
)

// ---------------------------------------------------------------------------
// HTML templates — direct conversion from .NET HtmlMessageTemplate.cs.
// Placeholders are mapped from C# {placeholder} to Go {{.Field}}.
// ---------------------------------------------------------------------------

const newCompetencyMessageToApprover = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that a new competency entry {{.CompetencyInfo}} has been created by {{.Requester}}, and your valuable review and approval is required.<br/><br/>
To review and act on this new competency entry, please follow the link provided below:
<strong><a href="{{.AppURL}}">Link&gt;&gt;</a> </strong> <br/><br/>
Your prompt attention to this matter is greatly appreciated.<br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const newCompetencyMessageToRequestor = `Dear <strong>{{.Recipient}}</strong>,
<br/><br/>
This is to notify you that a New Competency has been created and sent to your Head of Office for approval. Kindly log in <a href="{{.AppURL}}">here</a> for more details.
<br/><br/>
Regards`

const approvedCompetencyMessageToRequestor = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that the new Competencies Entries {{.CompetencyInfo}} has been Approved Successfuly. Kindly log in <a href="{{.AppURL}}">here</a> for more details.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const approvedCompetencyMessageToApprover = `Dear <strong>{{.Recipient}}</strong>,
<br/><br/>
This is to notify you that the New Competencies has been Approved Successfuly. Kindly log in <a href="{{.AppURL}}">here</a> for more details.
<br/><br/>
Regards`

const newReviewPeriodMessageToRequestor = `Dear <strong>{{.Recipient}}</strong>,
<br/><br/>
This is to notify you that your request to create new Review Period has been created and sent to your line manager for approval. Kindly log in <a href="{{.AppURL}}">here</a> for more details.
<br/><br/>
Regards`

const newReviewPeriodMessageToApprover = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that a new competency review period <strong>{{.ReviewPeriod}}</strong> has been created by <strong>{{.Requester}}</strong> and your valuable review and approval is required.<br/><br/>
Please click on the link below to review the created competency review period and take necessary action:<br/><br/>
<strong><a href="{{.AppURL}}">Link&gt;&gt;</a> </strong> <br/><br/>
Your prompt attention to this matter is greatly appreciated.<br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const approvedReviewPeriodMessageToAllStaff = `Dear Staff,<br/><br/>
The Human Resource Department has initiated a new competency review request for <strong>{{.ReviewPeriod}}</strong>. Kindly log in to the Competency Management System to thoroughly assess and objectively assign proficiency ratings to the competencies relevant to your Job and those of your colleagues. Your ratings should be based on an unbiased evaluation of how you and/or others have demonstrated mastery in the specified competencies up to the present date.<br/><br/>
To begin the competency review process, please click on the link provided below to access the Competency Management System.<br/><br/>
<strong><a href="{{.AppURL}}">Link&gt;&gt;</a></strong> <br/><br/>
Note that your competency review is considered incomplete until proficiency ratings are provided for all competencies assigned to you. This includes conducting a self-assessment as well as reviewing the competencies of other assigned parties, such as peers, superiors, or subordinates.<br/><br/>
Don't forget to click the submit button after completing the competency reviews for both you and other assigned parties.<br/><br/>
Your prompt response is much appreciated.<br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const reminderMessage = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is a gentle reminder to complete your <strong>{{.ReviewPeriod}}</strong> before the deadline of <strong>{{.DeadlineDate}}</strong>. Please note that you will not be able to submit your reviews after the deadline.<br/><br/>
To complete the competency review process, please click on the link provided below to access the Competency Management System.<br/><br/>
<strong><a href="{{.AppURL}}">Link&gt;&gt;</a></strong> <br/><br/>
Kindly make sure you complete the competency review for other assigned parties if you haven't done so. You may wish to note that, your competency review is considered incomplete until proficiency ratings are provided for all competencies assigned to you. This includes conducting a self-assessment as well as reviewing the competencies of other assigned parties, such as peers, superiors, or subordinates.<br/><br/>
Don't forget to click the <strong>Save Review</strong> button after completing the competency reviews for both you and other assigned parties.<br/><br/>
Thank you.<br/><br/>
Sincerely,<br/><br/>
<strong>Head, Talent Management Division.</strong><br/><br/>`

const selfReviewsCompletedMessageToStaff = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that you have completed the <strong>{{.ReviewPeriod}}</strong> competency review for yourself. Please ensure you complete the competency review for other assigned parties if you have not done so. <br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const peersReviewsCompletedMessageToStaff = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that the <strong>{{.ReviewPeriod}}</strong> competency review has been completed by your <strong>Peer</strong>.<br/><br/>
Please ensure you also complete the competency review for other assigned parties if you have not done so. <br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const subordinateReviewsCompletedMessageToSupervisor = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that the <strong>{{.ReviewPeriod}}</strong> competency review has been completed by your <strong>Subordinate</strong>.<br/><br/>
Please ensure you also complete the competency review for other assigned parties if you have not done so. <br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const superiorReviewsCompletedMessageToSubordinate = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that the <strong>{{.ReviewPeriod}}</strong> Competency review has been completed by your <strong>Superior</strong>.<br/><br/>
Please ensure you also complete the competency review for other assigned parties if you have not done so. <br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const developmentTaskAssignedMessageToStaff = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that you have been assigned a developmental task to help you improve your proficiency in a specific competency <strong>{{.CompetencyName}}</strong>.<br/><br/>
Please click on the link provided below to access the Competency Management System and review the assigned development task.<br/><br/>
<a href="{{.AppURL}}">Link&gt;&gt;</a> <br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const developmentTaskCompletedMessageToLineManager = `Dear <strong>{{.Recipient}}</strong>,
<br/><br/>
This is to notify you that a Development Task has been Completed by <strong>{{.RequestorName}}</strong>. Kindly log in <a href="{{.AppURL}}">here</a> for more details.<br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const developmentTaskApprovedMessageToStaff = `Dear <strong>{{.Recipient}}</strong>,
<br/><br/>
This is to notify you that a Development Task has been Approved by your Line Manager. Kindly log in <a href="{{.AppURL}}">here</a> for more details.<br/><br/> Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const jobRoleUpdateRequestMessageToApprover = `Dear <strong>{{.Recipient}}</strong>,
<br/><br/>
This is to notify you that there is a request for Job Role Update from <span style="font-weight:bold;">{{.Requester}}</span> requiring your attention.<br/><br/>
Please click on the link provided below to access the Competency Management System to Approve/Reject the request.<br/><br/>
<a href="{{.AppURL}}">Link&gt;&gt;</a> <br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const approvedJobRoleUpdateRequestMessage = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that the request to update your Job Role <strong>{{.Description}}</strong> has been approved. Kindly log in <a href="{{.AppURL}}">here</a> for more details.<br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const rejectJobRoleUpdateRequestMessage = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that the request to update your Job Role <strong>{{.Description}}</strong>.<br/><br/>
Kindly log in <a href="{{.AppURL}}">here</a> for more details.<br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

const cmsTeamApprovedJobRoleUpdateMessage = `Dear <strong>{{.Recipient}}</strong>,<br/><br/>
This is to notify you that the request to update your Job Role <strong>{{.Description}}</strong>.<br/><br/>
Thank you for your cooperation.<br/><br/>
Sincerely,<br/><br/>
<strong>CBN Competency Management System.</strong><br/><br/>`

// ---------------------------------------------------------------------------
// Template registry
// ---------------------------------------------------------------------------

var (
	templateMap  map[string]*template.Template
	templateOnce sync.Once
)

func initTemplates() {
	templateMap = make(map[string]*template.Template, 19)

	mustParse := func(key, body string) {
		t, err := template.New(key).Parse(body)
		if err != nil {
			panic(fmt.Sprintf("email: failed to parse template %s: %v", key, err))
		}
		templateMap[key] = t
	}

	mustParse(TplNewCompetencyToApprover, newCompetencyMessageToApprover)
	mustParse(TplNewCompetencyToRequestor, newCompetencyMessageToRequestor)
	mustParse(TplApprovedCompetencyToRequestor, approvedCompetencyMessageToRequestor)
	mustParse(TplApprovedCompetencyToApprover, approvedCompetencyMessageToApprover)
	mustParse(TplNewReviewPeriodToRequestor, newReviewPeriodMessageToRequestor)
	mustParse(TplNewReviewPeriodToApprover, newReviewPeriodMessageToApprover)
	mustParse(TplApprovedReviewPeriodToAll, approvedReviewPeriodMessageToAllStaff)
	mustParse(TplReminderMessage, reminderMessage)
	mustParse(TplSelfReviewCompleted, selfReviewsCompletedMessageToStaff)
	mustParse(TplPeersReviewCompleted, peersReviewsCompletedMessageToStaff)
	mustParse(TplSubordinateReviewCompleted, subordinateReviewsCompletedMessageToSupervisor)
	mustParse(TplSuperiorReviewCompleted, superiorReviewsCompletedMessageToSubordinate)
	mustParse(TplDevTaskAssigned, developmentTaskAssignedMessageToStaff)
	mustParse(TplDevTaskCompleted, developmentTaskCompletedMessageToLineManager)
	mustParse(TplDevTaskApproved, developmentTaskApprovedMessageToStaff)
	mustParse(TplJobRoleUpdateRequest, jobRoleUpdateRequestMessageToApprover)
	mustParse(TplApprovedJobRole, approvedJobRoleUpdateRequestMessage)
	mustParse(TplRejectJobRole, rejectJobRoleUpdateRequestMessage)
	mustParse(TplCMSTeamApprovedJobRole, cmsTeamApprovedJobRoleUpdateMessage)
}

// Render executes the named template with the given data and returns the
// rendered HTML string. Returns an error if the template key is unknown or
// execution fails.
func Render(templateKey string, data TemplateData) (string, error) {
	templateOnce.Do(initTemplates)

	t, ok := templateMap[templateKey]
	if !ok {
		return "", fmt.Errorf("email: unknown template key %q", templateKey)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("email: rendering template %s: %w", templateKey, err)
	}
	return buf.String(), nil
}
