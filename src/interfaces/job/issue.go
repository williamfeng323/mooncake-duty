package jobinterface

import (
	"github.com/robfig/cron/v3"
)

// RegisterIssueJobs assign the job to the inputted job instance
func RegisterIssueJobs(c *cron.Cron) {
	c.AddFunc("@every 1m", notifyIssues)
}

func notifyIssues() {
	/*
		issues := issueService.GetIssueLists()
		for issue of each issue,
			if issue.shouldNotify()
				go issue.NotifyAsync(completeChan)
	*/
	// issueService := issue.GetIssueService()
}
