package routers

import (
	"BubbleList/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 告诉gin框架模板文件引用的静态文件在哪
	r.Static("/static", "src/static")
	// 告诉gin框架模板文件在哪
	r.LoadHTMLGlob("src/templates/*")
	r.GET("/", controller.IndexHandler)
	// 定义路由组
	v1Group := r.Group("v1")
	{
		// 待办事项
		// 添加待办事情
		v1Group.POST("/todo", controller.CreateATodo)
		// 查看所有待办事项
		v1Group.GET("/todo", controller.GetTodoList)
		// 更新某一待办事项的完成状态
		v1Group.PUT("/todo/:id", controller.UpdateATodo)
		// 删除某一待办事情
		v1Group.DELETE("/todo/:id", controller.DeleteATodo)
	}
	return r
}
