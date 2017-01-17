package entity

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"time"
)

// JobInfo  job信息详情

type JobInfo struct {
	Id                  int `orm:"pk;auto"`
	Name                string
	Group               string
	Type                string
	Time                time.Time
	Cron                string
	Urls                string
	ClassPath           string
	InvokePolicy        string
	IsActivity          int
	Desc                string
	CreateTime          time.Time
	ModifyTime          time.Time
	Param               string
	LatestTriggerTime   time.Time
	LatestServerAddress string
	OwnerPhone          string
}

func (this *JobInfo) FindAllJobInfo() ([]*JobInfo, error) {

	var jobs []*JobInfo
	o := orm.NewOrm()
	qs := o.QueryTable("job_info")
	qs = qs.Filter("is_activity", 1)
	_, err := qs.OrderBy("id", "-modify_time", "-create_time").All(&jobs)

	//common.PanicIf(err)
	return jobs, err

}

func (this *JobInfo) FindAllJobInfoByPage() ([]*JobInfo, error) {

	var jobs []*JobInfo
	o := orm.NewOrm()
	qs := o.QueryTable("job_info")
	if this.Name != "" {
		qs = qs.Filter("name", this.Name)
	}

	if this.Group != "" {
		qs = qs.Filter("group", this.Group)
	}
	if this.Name == "" && this.Group == "" {
		qs = qs.Limit(100)
	}

	_, err := qs.OrderBy("id", "-modify_time", "-create_time").All(&jobs)

	//common.PanicIf(err)
	return jobs, err

}

func (this *JobInfo) SaveJobInfo() error {

	_, err := orm.NewOrm().Insert(this)

	return err
}

func (this *JobInfo) GetJobInfoById() error {

	o := orm.NewOrm()
	return o.Read(this)
}

func (this *JobInfo) UpdateJobInfo() error {

	o := orm.NewOrm()
	id, err := o.Update(this, "name", "group", "type", "time", "cron", "urls", "class_path", "invoke_policy", "is_activity", "desc", "modify_time", "param", "owner_phone")
	if err != nil {
		return err
	}

	if id == 0 {
		return errors.New("此记录不存在!")
	}
	return nil
}

func (this *JobInfo) ActiveJobInfo() error {

	o := orm.NewOrm()
	id, err := o.Update(this, "is_activity", "modify_time")
	if err != nil {
		return err
	}

	if id == 0 {
		return errors.New("此记录不存在!")
	}
	return nil
}

func (this *JobInfo) GetJobInfo() error {

	o := orm.NewOrm()
	return o.Read(this)

}

func (this *JobInfo) DeleteJobInfo() error {

	err := this.GetJobInfo()
	o := orm.NewOrm()
	defer func(err error) {

		if err != nil {
			o.Rollback()
		} else {
			o.Commit()
		}
	}(err)

	if err == nil {

		_, err = o.Delete(this)
		if err == nil {
			jobInfoHistory := &JobInfoHistory{
				Id:                  this.Id,
				Name:                this.Name,
				Group:               this.Group,
				Type:                this.Type,
				Time:                this.Time,
				Cron:                this.Cron,
				Urls:                this.Urls,
				ClassPath:           this.ClassPath,
				InvokePolicy:        this.InvokePolicy,
				IsActivity:          this.IsActivity,
				Desc:                this.Desc,
				CreateTime:          this.CreateTime,
				ModifyTime:          this.ModifyTime,
				Param:               this.Param,
				LatestTriggerTime:   this.LatestTriggerTime,
				LatestServerAddress: this.LatestServerAddress,
				OwnerPhone:          this.OwnerPhone,
			}

			err = jobInfoHistory.SaveJobInfoHistory()

		}

	}

	return err
}
