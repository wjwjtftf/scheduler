package entity

import (
	"github.com/astaxie/beego/orm"
	"log"
	"time"
)

const (
	INIT      = "INIT"
	INVOKING  = "INVOKING"
	EXECUTING = "EXECUTING"
	COMPLETED = "COMPLETED"
	ERROR     = "ERROR"
	FINISH
)

// job 快照
type JobSnapshot struct {
	Id            int `orm:"pk;auto"`
	JobInfoId     int
	Name          string
	Group         string
	Status        string
	Ip            string
	Url           string
	Result        string
	TimeConsume   int64
	Detail        string
	CreateTime    time.Time
	ModifyTime    time.Time
	ServerAddress string
	Params        string
	NextTime      time.Time
}

// find jobSnapshot list max 1000
func (this *JobSnapshot) FindList() ([]*JobSnapshot, error) {

	var jobSnapshotList []*JobSnapshot
	o := orm.NewOrm()
	qs := o.QueryTable("job_snapshot")
	if this.Name != "" {
		qs = qs.Filter("name", this.Name)
	}

	if this.Group != "" {
		qs = qs.Filter("group", this.Group)
	}

	if this.Status != "" {
		qs = qs.Filter("status", this.Status)
	}
	if this.Status == "" && this.Group == "" && this.Name == "" {
		qs = qs.Limit(100)
	}
	_, err := qs.OrderBy("-modify_time", "-create_time").All(&jobSnapshotList)
	return jobSnapshotList, err

}

func (this *JobSnapshot) GetJobSnapshot() error {

	o := orm.NewOrm()
	return o.Read(this)
}

func (this *JobSnapshot) InsertJobSnapshot() error {
	o := orm.NewOrm()
	_, err := o.Insert(this)
	return err
}

func (this *JobSnapshot) UpdateSnapshot() error {

	log.Println("modifyTime= ", this.ModifyTime)
	o := orm.NewOrm()
	_, err := o.Update(this, "status", "ip", "url", "result", "time_consume", "detail", "modify_time", "server_address", "params")

	return err
}
