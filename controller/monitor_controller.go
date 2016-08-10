package controller

import (
	"github.com/wjwjtftf/scheduler/common"
	"github.com/wjwjtftf/scheduler/job"
)

type MonitorController struct {
	BaseController
}

func (this *MonitorController) Index() {

	jobManger := job.NewJobMnager()
	jobList, err := jobManger.GetJobSnapshotList()
	common.PanicIf(err)

	this.TplName = "monitor/index.html"
	this.Data["jobList"] = jobList
	this.Render()
}
