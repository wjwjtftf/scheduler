package quartz

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

var (
	Error_EXISTS_JOB     = errors.New("exists job error !")
	Error_NOT_EXISTS_JOB = errors.New("not exists job error !")
)

// defind the job call method
type JobFunc func(jobId int, targetUrl, params string, nextTime time.Time)

// The Schedule describes a job's duty cycle.
type Schedule interface {
	// Return the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next(time.Time) time.Time
}

// job  多余字段为了将来进行扩展
type Job struct {
	// id
	Id int
	// 任务名称
	Name string
	// 分组
	Group string
	// ip
	Ip string
	//目标服务器地址
	Url string
	// 上次触发时间
	Pre time.Time
	// 下次执行时间
	Next time.Time
	//schedule
	schedule Schedule
	// 执行的job 方法
	JobFunc JobFunc
	// corn 表达式
	Expression string
	// 状态 0 init  1 invoking 2 COMPLETED 3
	State uint
	//请求body字符串
	Params string
	//执行策略
	InvokePolicy string
	//是否激活  1 激活 0 非激活
	Active int
}

type Quartz struct {
	jobPool   []*Job
	Stop      chan int
	addJob    chan *Job
	removeJob chan int
	running   bool
	snapshot  chan []*Job
	modifyJob chan *Job
	jobError  chan error
	lk        sync.Mutex
}

func New() *Quartz {

	return &Quartz{
		Stop:      make(chan int),
		addJob:    make(chan *Job),
		removeJob: make(chan int),
		modifyJob: make(chan *Job),
		snapshot:  make(chan []*Job),
		jobError:  make(chan error),
	}
}

// add a job  to the
func (qz *Quartz) AddJob(job *Job) error {
	qz.lk.Lock()
	defer qz.lk.Unlock()
	if qz.hasExistsJob(job.Id) == -1 {
		log.Println(job.Expression)
		schedule, err := Parse(job.Expression)

		if err != nil {
			return err
		}
		return qz.scheduleJob(schedule, job)
	} else {
		return Error_EXISTS_JOB
	}

}

// set schedule to job
func (qz *Quartz) scheduleJob(schedule Schedule, job *Job) error {
	job.schedule = schedule
	if qz.running == false {
		qz.jobPool = append(qz.jobPool, job)
		return nil
	} else {
		qz.addJob <- job
	}

	return <-qz.jobError

}

// if has exists the same job
func (qz *Quartz) hasExistsJob(JobId int) int {

	for i, v := range qz.jobPool {

		if v.Id == JobId {
			return i
		}
	}
	return -1
}

func (qz *Quartz) BootStrap() {

	go qz.run()

}

// stop the quartz
func (qz *Quartz) StopJob() {
	qz.lk.Lock()
	defer qz.lk.Unlock()
	qz.Stop <- 1

}

// remove a job by jobName
func (qz *Quartz) RemoveJob(jobId int) error {
	qz.lk.Lock()
	defer qz.lk.Unlock()
	index := qz.hasExistsJob(jobId)
	if index == -1 {
		return Error_NOT_EXISTS_JOB
	}

	if !qz.running {
		qz.jobPool = qz.jobPool[:index+copy(qz.jobPool[index:], qz.jobPool[index+1:])]
		return nil
	} else {
		qz.removeJob <- jobId
	}

	return <-qz.jobError
}

// modify jobinfo
func (qz *Quartz) ModifyJob(job *Job) error {
	qz.lk.Lock()
	defer qz.lk.Unlock()
	index := qz.hasExistsJob(job.Id)
	if index == -1 {
		return Error_NOT_EXISTS_JOB
	}
	schedule, err := Parse(job.Expression)
	if err != nil {
		return err
	}
	job.schedule = schedule
	if !qz.running {
		qz.jobPool[index] = job
		return nil

	} else {
		log.Println("modify job....job = ", job)
		qz.modifyJob <- job

		return <-qz.jobError
	}

}
func (qz *Quartz) SnapshotJob() ([]*Job, error) {

	if !qz.running {

		return qz.jobPool, nil
	} else {

		/*
			qz.snapshot <- nil
			list := <-qz.snapshot
			return list,nil
		*/

		return qz.jobPool, nil
	}
}

// quartz run
func (qz *Quartz) run() {
	fmt.Println("quartz is run.....")
	qz.running = true
	now := time.Now().Local()

	for _, v := range qz.jobPool {
		v.Next = v.schedule.Next(now)
	}

	for {

		sort.Sort(JobPool(qz.jobPool))

		var effective time.Time
		if len(qz.jobPool) == 0 || qz.jobPool[0].Next.IsZero() {
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = qz.jobPool[0].Next
		}

		select {

		case now = <-time.After(effective.Sub(now)):

			for _, v := range qz.jobPool {

				if v.Next != effective {
					break
				}
				v.Pre = effective
				log.Println("effective is ", effective, "next time is ", v.schedule.Next(effective))
				v.Next = v.schedule.Next(effective)
				go v.JobFunc(v.Id, v.Url, v.Params, v.Next)
				// 根据策略删除一次性任务
				if v.InvokePolicy == "once" {
					qz.removeJob <- v.Id
				}
			}
			continue
		// add new job
		case newJob := <-qz.addJob:
			log.Println("recev the add job")
			if qz.hasExistsJob(newJob.Id) == -1 {
				now = time.Now().Local()
				qz.jobPool = append(qz.jobPool, newJob)
				fmt.Println("now is ", now, "add effective is ", newJob.schedule.Next(now))
				newJob.Next = newJob.schedule.Next(now)
				qz.jobError <- nil

			} else {
				qz.jobError <- Error_EXISTS_JOB
			}
		// remove the exists  job
		case removeJobId := <-qz.removeJob:

			index := qz.hasExistsJob(removeJobId)
			log.Println("remove a job removeJobId= ", removeJobId)
			if index != -1 {
				qz.jobPool = qz.jobPool[:index+copy(qz.jobPool[index:], qz.jobPool[index+1:])]
				qz.jobError <- nil
			} else {
				qz.jobError <- Error_NOT_EXISTS_JOB
			}

		case <-qz.Stop:
			log.Println("stop the quartz!")
			qz.running = false
			return

		case <-qz.snapshot:
			qz.snapshot <- qz.jobPool

		case modifyJob := <-qz.modifyJob:
			log.Println("modify a job ", modifyJob)
			index := qz.hasExistsJob(modifyJob.Id)
			if index == -1 {
				log.Println("=#####################=")
				qz.jobError <- Error_NOT_EXISTS_JOB
			} else {
				now = time.Now().Local()
				nextTime := modifyJob.schedule.Next(now)
				fmt.Println("modify  job now is ", now, "add effective is ")
				modifyJob.Next = nextTime
				qz.jobPool[index] = modifyJob
				log.Println("=***************=", modifyJob, "nextTime = ", nextTime)
				qz.jobError <- nil
			}

		case err := <-qz.jobError:
			if err != nil {
				log.Println(err)
			}
		}

	}

}

type JobPool []*Job

// Len is the number of elements in the collection.
func (jp JobPool) Len() int {
	return len(jp)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (jp JobPool) Less(i, j int) bool {

	if jp[i].Next.IsZero() {
		return false
	}

	if jp[j].Next.IsZero() {
		return true
	}

	return jp[i].Next.Before(jp[j].Next)
}

// Swap swaps the elements with indexes i and j.
func (jp JobPool) Swap(i, j int) {
	var temp = jp[i]
	jp[i] = jp[j]
	jp[j] = temp
}
