package controller

type AboutController struct  {
	BaseController

}

func (this *AboutController)Index() {

	this.TplName = "about/about.html"
	this.Render()
}