package controller

type IndexController  struct{
	BaseController

}

func (this *IndexController) Index()  {

	this.TplName = "index.html"
	this.Render()
}

func (this *IndexController) User()  {

	this.TplName = "user/user.html"
	this.Render()

}

