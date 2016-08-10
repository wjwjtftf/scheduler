package invoke

import (
	"bytes"
	"encoding/json"
	"github.com/wjwjtftf/scheduler/common"
	"github.com/wjwjtftf/scheduler/entity"
	"github.com/wjwjtftf/scheduler/policy"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// 执行
type Invoker struct {
}

// 执行任务
func (this *Invoker) Execute(jobInfo *entity.JobInfo, nextTime time.Time, params string) error {

	snapshot, err := this.Init(jobInfo, nextTime)
	if err != nil {
		return err
	}

	log.Println("snapshot = ", snapshot)

	fac := &policy.Factory{}
	policy := fac.FindPolicy(jobInfo)

	for {
		url := policy.GetNextUrl()

		if url == "" {
			//this.executeJobResult(snapshot, time.Now().Local().Format("2006-01-02 15:04:05")+"所有目标服务器地址都不可用", entity.ERROR)
			break
		}
		snapshot.Url = url

		// 准备执行
		err = this.invoke(snapshot)

		if err == nil {
			break
		}
	}

	return err
}

// 执行任务
func (this *Invoker) invoke(jobSnapshot *entity.JobSnapshot) error {

	this.executeJobResult(jobSnapshot, common.Now()+"准备任务提交至目标服务器地址:"+jobSnapshot.Url, entity.INVOKING)

	//startTime := time.Now()

	//jobRequest := &common.JobRequest{JobSnapshot: jobSnapshot.Id, Params: jobSnapshot.Params, Status: entity.INVOKING}
	//content, err := json.Marshal(jobRequest)

	content := jobSnapshot.Params
	if content == "" {
		content = "{}"
	}

	resp, err := http.Post(jobSnapshot.Url, "application/json;charset=utf-8", bytes.NewBuffer([]byte(content)))
	if err != nil {
		this.executeJobResult(jobSnapshot, common.Now()+"目标服务器地址:"+jobSnapshot.Url+"不可用", entity.ERROR)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	result := &common.JobResponse{}
	log.Println("body = ", string(body))

	err = json.Unmarshal(body, result)
	if err != nil {
		this.executeJobResult(jobSnapshot, common.Now()+"目标服务器地址:"+jobSnapshot.Url+"非法的响应["+string(body)+"]", entity.ERROR)
		log.Println("err = ", err)
		return err
	}

	addr := resp.Request.RemoteAddr

	log.Println("serveraddr=", addr)
	jobSnapshot.ServerAddress = addr
	if result.Success == true {
		//nowTime := time.Now()
		//d := nowTime.Sub(startTime)
		//timeConsume := int64(d.Seconds())
		jobSnapshot.TimeConsume = result.ResponseTime
		jobSnapshot.Result = string(body)
		jobSnapshot.Ip = common.GetIPFromUrl(jobSnapshot.Url)
		this.executeJobResult(jobSnapshot, common.Now()+"目标服务器地址:"+jobSnapshot.Url+"执行任务已经成功完成", entity.COMPLETED)
		//如果是异步job，则需要选择下面的方式不断检查目标服务器的执行状态
		//go this.processCheckJobResult(jobSnapshot)

	} else {
		this.executeJobResult(jobSnapshot, common.Now()+"目标服务器地址:"+jobSnapshot.Url+"执行失败["+result.Error.Message+"]", entity.ERROR)
	}

	return nil
}

func (this *Invoker) executeJobResult(snapshot *entity.JobSnapshot, detail, status string) {

	snapshot.Detail = snapshot.Detail + detail + "\n"
	snapshot.Status = status
	snapshot.ModifyTime = time.Now().Local()
	snapshot.UpdateSnapshot()
}

func (this *Invoker) Init(jobInfo *entity.JobInfo, nextTime time.Time) (*entity.JobSnapshot, error) {

	serverAddr := common.GetLocalAddr()
	snapshot := &entity.JobSnapshot{
		JobInfoId:     jobInfo.Id,
		Name:          jobInfo.Name,
		Group:         jobInfo.Group,
		Status:        entity.INIT,
		Url:           jobInfo.Urls,
		TimeConsume:   0,
		ServerAddress: serverAddr,
		NextTime:      nextTime,
		CreateTime:    time.Now().Local(),
		Detail:        common.Now() + "初始化" + "\n",
		Params:        jobInfo.Param,
	}
	err := snapshot.InsertJobSnapshot()
	return snapshot, err
}

// 执行更新状态
func (this *Invoker) processCheckJobResult(jobSnapshot *entity.JobSnapshot) {
	var quit bool = false
	for !quit {

		select {

		case <-time.After(time.Second * 5):

			jobRequest := &common.JobRequest{JobSnapshot: jobSnapshot.Id, Params: jobSnapshot.Params, Status: jobSnapshot.Status}

			content, err := json.Marshal(jobRequest.Params)
			if err != nil {
				this.executeJobResult(jobSnapshot, common.Now()+"解析job请求参数出错", entity.ERROR)
				continue
			}

			resp, err := http.Post(jobSnapshot.Url, "application/json;charset=utf-8", bytes.NewBuffer(content))
			if err != nil {
				this.executeJobResult(jobSnapshot, common.Now()+"目标服务器地址:"+jobSnapshot.Url+"不可用", entity.ERROR)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				resp.Body.Close()
				continue
			}
			result := &common.JobResponse{}
			log.Println("body = ", string(body))
			err = json.Unmarshal(body, result)
			if err != nil {
				resp.Body.Close()
				continue
			}
			log.Println("result= ", result)
			if result.Status == entity.EXECUTING {
				this.executeJobResult(jobSnapshot, common.Now()+"目标服务器地址:"+jobSnapshot.Url+" 正在执行中...", entity.EXECUTING)
			} else if result.Status == entity.COMPLETED {
				this.executeJobResult(jobSnapshot, common.Now()+"目标服务器地址:"+jobSnapshot.Url+"任务执行完成...", entity.COMPLETED)
				quit = true
			} else {
				this.executeJobResult(jobSnapshot, common.Now()+"目标服务器地址:"+jobSnapshot.Url+"任务执行失败...", entity.ERROR)
				quit = true
			}
		}
	}
}
