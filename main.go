package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wjwjtftf/scheduler/controller"
	"github.com/wjwjtftf/scheduler/entity"
	"github.com/wjwjtftf/scheduler/job"
	"runtime"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	//"scheduler:scheduler@tcp(127.0.0.1:13306)/scheduler?charset=utf8&loc=Local"
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&loc=Local"
	orm.RegisterDataBase("default", "mysql", dsn)

	orm.RegisterModel(&entity.JobInfo{}, &entity.JobInfoHistory{}, &entity.JobSnapshot{})
}

func main() {

	// set CPU
	runtime.GOMAXPROCS(runtime.NumCPU())
	jobManager := job.NewJobMnager()
	jobManager.PushAllJob()

	if beego.AppConfig.String("runmode") == "dev" {
		//orm.Debug = true
		beego.SetLevel(beego.LevelDebug)
	} else {
		beego.SetLevel(beego.LevelInformational)
	}

	// TODO Init jobList

	// set home  path
	beego.Router("/", &controller.IndexController{}, "get:Index")

	// jobinfo
	beego.Router("/jobinfo/list", &controller.JobInfoManagerController{}, "*:List")
	beego.Router("/jobinfo/add", &controller.JobInfoManagerController{}, "get:ToAdd")
	beego.Router("/jobinfo/add", &controller.JobInfoManagerController{}, "post:Add")
	beego.Router("/jobinfo/edit", &controller.JobInfoManagerController{}, "get:ToEdit")
	beego.Router("/jobinfo/edit", &controller.JobInfoManagerController{}, "post:Edit")
	beego.Router("/jobinfo/delete", &controller.JobInfoManagerController{}, "post:Delete")
	beego.Router("/jobinfo/info", &controller.JobInfoManagerController{}, "get:Info")
	beego.Router("/jobinfo/active", &controller.JobInfoManagerController{}, "*:Active")

	// jobsnapshot
	beego.Router("/jobsnapshot/list", &controller.JobSnapshotController{}, "*:List")
	beego.Router("/jobsnapshot/info", &controller.JobSnapshotController{}, "get:Info")

	// jobinfohistory
	beego.Router("/jobinfohistory/list", &controller.JobInfoHistoryController{}, "*:List")

	//about
	beego.Router("/about", &controller.AboutController{}, "*:Index")

	//monitor
	beego.Router("/monitor/", &controller.MonitorController{}, "*:Index")

	// set static resource
	beego.SetStaticPath("static", "static")
	beego.SetStaticPath("public", "static")

	// start web app
	beego.Run()

}
