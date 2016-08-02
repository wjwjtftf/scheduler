package controller

import (
	"scheduler/entity"
	"scheduler/common"
)

// 任务详情历史控制器
type JobInfoHistoryController struct  {

	 BaseController

}

// jobinfohistory list
func  (this *JobInfoHistoryController) List()  {

	name := this.GetString("Name")
	group := this.GetString("Group")

	history := entity.JobInfoHistory{Name:name,Group:group}
	historys,err := history.FindAllJobInfoList()

	common.PanicIf(err)
	this.Data["historys"] = historys
	this.Data["name"] = name
	this.Data["group"] = group
	this.TplName = "jobinfohistory/list.html"
	this.Render()

}
