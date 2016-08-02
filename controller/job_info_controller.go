package controller

import (
	"scheduler/common"
	"scheduler/entity"
	"scheduler/job"
	"time"
)

// job info
type JobInfoManagerController struct {
	BaseController
}

func (this *JobInfoManagerController) List() {
	name := this.GetString("Name")
	group := this.GetString("Group")
	jobInfo := &entity.JobInfo{Name: name, Group: group}
	jobs, err := jobInfo.FindAllJobInfoByPage()
	common.PanicIf(err)
	this.Data["jobs"] = jobs
	this.Data["group"] = group
	this.Data["name"] = name
	this.TplName = "jobinfo/list.html"
	this.Render()
}

func (this *JobInfoManagerController) ToAdd() {

	this.TplName = "jobinfo/add.html"
	this.Render()
}

func (this *JobInfoManagerController) Add() {

	jsonResult := &entity.JsonResult{}

	name := this.GetString("Name")
	group := this.GetString("Group")
	cron := this.GetString("Cron")
	urls := this.GetString("Urls")
	invokePolicy := this.GetString("InvokePolicy")
	param := this.GetString("Param")
	desc := this.GetString("Desc")
	ownerPhone := this.GetString("OwnerPhone")

	if name == "" {
		jsonResult.Message = "job名称不能为空"
	} else if group == "" {
		jsonResult.Message = "job任务分组不能为空"
	} else if invokePolicy == "" {
		jsonResult.Message = "执行策略不能为空"
	} else if urls == "" {
		jsonResult.Message = "目标服务器urls不能为空"
	} else if cron == "" {
		jsonResult.Message = "Cron表达式不能为空"
	}
	if jsonResult.Message == "" {
		jobInfo := &entity.JobInfo{}
		jobInfo.CreateTime = time.Now()
		jobInfo.Cron = cron
		jobInfo.Desc = desc
		jobInfo.Group = group
		jobInfo.InvokePolicy = invokePolicy
		jobInfo.Param = param
		jobInfo.Name = name
		jobInfo.Urls = urls
		jobInfo.OwnerPhone = ownerPhone
		jobInfo.ModifyTime = time.Now()
		err := jobInfo.SaveJobInfo()
		if err != nil {

			jsonResult.Message = "保存失败,请重试!"
		} else {

			jsonResult.Message = "保存成功"
			jsonResult.Success = true

		}
	}

	this.Data["json"] = jsonResult
	this.ServeJSON()
}

func (this *JobInfoManagerController) ToEdit() {

	id, err := this.GetInt("Id")
	common.PanicIf(err)

	jobInfo := &entity.JobInfo{Id: id}
	err = jobInfo.GetJobInfoById()
	common.PanicIf(err)
	this.Data["jobInfo"] = jobInfo
	this.TplName = "jobinfo/edit.html"
	this.Render()
}

func (this *JobInfoManagerController) Info() {

	id, err := this.GetInt("Id")
	common.PanicIf(err)

	jobInfo := &entity.JobInfo{Id: id}
	err = jobInfo.GetJobInfoById()
	common.PanicIf(err)
	this.Data["jobInfo"] = jobInfo
	this.TplName = "jobinfo/info.html"
	this.Render()
}
func (this *JobInfoManagerController) Edit() {
	jsonResult := &entity.JsonResult{}
	id, err := this.GetInt("Id")
	name := this.GetString("Name")
	group := this.GetString("Group")
	cron := this.GetString("Cron")
	urls := this.GetString("Urls")
	invokePolicy := this.GetString("InvokePolicy")
	param := this.GetString("Param")
	desc := this.GetString("Desc")
	ownerPhone := this.GetString("OwnerPhone")

	if name == "" {
		jsonResult.Message = "job名称不能为空"
	} else if err != nil {
		jsonResult.Message = "此记录不存在!"
	} else if group == "" {
		jsonResult.Message = "job任务分组不能为空"
	} else if invokePolicy == "" {
		jsonResult.Message = "执行策略不能为空"
	} else if urls == "" {
		jsonResult.Message = "目标服务器urls不能为空"
	} else if cron == "" {
		jsonResult.Message = "Cron表达式不能为空"

	}
	if jsonResult.Message == "" {
		jobInfo := &entity.JobInfo{}
		jobInfo.Id = id
		jobInfo.CreateTime = time.Now()
		jobInfo.Cron = cron
		jobInfo.Desc = desc
		jobInfo.Group = group
		jobInfo.InvokePolicy = invokePolicy
		jobInfo.Param = param
		jobInfo.Name = name
		jobInfo.Urls = urls
		jobInfo.OwnerPhone = ownerPhone
		jobInfo.ModifyTime = time.Now()
		err := jobInfo.UpdateJobInfo()
		if err != nil {

			jsonResult.Message = err.Error()
		} else {
			jsonResult.Message = "更新成功"
			jsonResult.Success = true
			jobManger := job.NewJobMnager()
			jobManger.DeleteJob(jobInfo.Id)
		}
	}

	this.Data["json"] = jsonResult
	this.ServeJSON()

}

// delete jobinfo

func (this *JobInfoManagerController) Delete() {
	jsonResult := &entity.JsonResult{}
	id, err := this.GetInt("Id")

	if err != nil {
		jsonResult.Message = "此记录不存在"
	} else {

		jobManger := job.NewJobMnager()
		jobManger.DeleteJob(id)
		jobInfo := entity.JobInfo{Id: id}
		err = jobInfo.DeleteJobInfo()
		if err != nil {
			jsonResult.Message = "删除失败,请重试"
		} else {
			jsonResult.Message = "删除成功"
			jsonResult.Success = true
		}

	}
	this.Data["json"] = jsonResult
	this.ServeJSON()
}

func (this *JobInfoManagerController) Active() {
	jsonResult := &entity.JsonResult{}
	id, err := this.GetInt("Id")
	active, err := this.GetInt("active")
	if err != nil {
		jsonResult.Message = "此记录不存在"
	} else {
		jsonResult.Data = active
		jobInfo := &entity.JobInfo{Id: id}
		err = jobInfo.GetJobInfoById()
		if err != nil {
			jsonResult.Message = "此job不存在"
		} else {

			jobInfo.IsActivity = active

			if active == 0 {
				jobManger := job.NewJobMnager()
				err = jobManger.DeleteJob(jobInfo.Id)

			} else if active == 1 {
				jobManger := job.NewJobMnager()
				err = jobManger.AddJob(jobInfo)
			}
			jobInfo.ModifyTime = time.Now()
			err = jobInfo.ActiveJobInfo()
			if err != nil {
				jsonResult.Message = "操作失败"
			} else {
				jsonResult.Message = "操作成功"
				jsonResult.Success = true
			}

		}

	}

	this.Data["json"] = jsonResult
	this.ServeJSON()
}
