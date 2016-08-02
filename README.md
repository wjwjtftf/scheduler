###  quartz of golang
###### 在企业项目开发中会定时执行对应的job，对于一些简单少的job可以直接使用调度器调度执行任务。当随着公司的业务越来越多，执行任务越来越多。那么直接使用任务调度器调度任务执行会变得臃肿，而且对于任务是动态配置不可实现。如：想某一个时刻停止任务的执行、删除此任务然后修改更新任务执行时间等，如某一个任务配置到多台机器上如何做到不可用时，进行转移等问题。

###### 项目参考

  1. https://github.com/shotdog/scheduler
  2. https://github.com/robfig/cron

#### 1、pre

* install golang env

  https://github.com/golang/go

#### 2、install  

 -       cd $GOPATH/src 
 -       git clone https://github.com/wjwjtftf/scheduler 
 -       go get github.com/astaxie/beego
 -       go get  github.com/go-sql-driver/mysql
 -       
 -       init db  scheduler.sql
 -       modify conf/app.conf -->db config


#### 3、run

*       cd $GOPATH
*       cd src
*       cd scheduler
*       go build main.go
*       ./main


#### 4、Screenshot

![image](https://github.com/wjwjtftf/scheduler/raw/master/screenshot/1.png)

![image](https://github.com/wjwjtftf/scheduler/raw/master/screenshot/2.png)

![image](https://github.com/wjwjtftf/scheduler/raw/master/screenshot/3.png)

![image](https://github.com/wjwjtftf/scheduler/raw/master/screenshot/4.png)

![image](https://github.com/wjwjtftf/scheduler/raw/master/screenshot/5.png)

![image](https://github.com/wjwjtftf/scheduler/raw/master/screenshot/6.png)

![image](https://github.com/wjwjtftf/scheduler/raw/master/screenshot/7.png)

#### 5、Protocol
* see [invoker.go](https://github.com/wjwjtftf/scheduler/blob/master/invoker/invoker.go)

#### 6、Client Test

* see [scheduler-client](https://github.com/wjwjtftf/scheduler-client-test)

   * cd scheduler-client
   * go build main.go
   * ./main
