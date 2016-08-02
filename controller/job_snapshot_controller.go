package controller

import (
	"scheduler/entity"
	"scheduler/common"
)

type JobSnapshotController struct {
	BaseController
}

func (this *JobSnapshotController) List() {

	name := this.GetString("Name");
	group := this.GetString("Group");
	status := this.GetString("Status")
	jobSnapshot := entity.JobSnapshot{Name:name, Group:group,Status:status}

	jobSnapshotList, err := jobSnapshot.FindList()
	common.PanicIf(err)
	this.Data["jobSnapshotList"] = jobSnapshotList
	this.Data["name"] = name
	this.Data["group"] = group
	this.Data["status"] = status

	this.TplName = "jobsnapshot/list.html"
	this.Render()
}

func (this *JobSnapshotController) ToAdd() {

	this.TplName = "jobsnapshot/add.html"
	this.Render()
}

func (this *JobSnapshotController) Add() {

	this.TplName = "jobsnapshot/add.html"
	this.Render()
}

func (this *JobSnapshotController) ToEdit() {

	this.TplName = "jobsnapshot/edit.html"
	this.Render()
}

func (this *JobSnapshotController) Edit() {

	this.TplName = "jobsnapshot/edit.html"
	this.Render()
}

func (this *JobSnapshotController) Info()  {

	  id,err:= this.GetInt("Id")
	common.PanicIf(err)

	jobSnapshot := entity.JobSnapshot{Id:id}

	err = jobSnapshot.GetJobSnapshot()
	common.PanicIf(err)
	this.Data["jobSnapshot"] = jobSnapshot
	this.TplName = "jobsnapshot/info.html"
	this.Render()
}