package controller

import (
	"BubbleList/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	url --> controller --> logic -->  models
请求来了 -->   控制器    --> 业务逻辑 --> 模型层CRUD
*/
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func CreateATodo(c *gin.Context) {
	// 前端页面填写后点击提交，请求发送到这里
	// 1. 从请求中得到数据
	var todo models.Todo
	_ = c.BindJSON(&todo)

	// 2. 存入数据库、返回响应
	err := models.CreateATodo(&todo)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, todo)
	}
}

func GetTodoList(c *gin.Context) {
	// 查询表里所有数据
	var todoList []models.Todo

	err := models.GetTodoList(&todoList)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, todoList)
	}
}

func UpdateATodo(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}

	todo, err := models.GetATodo(id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	_ = c.BindJSON(todo)

	err = models.UpdateATodo(todo)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, todo)
	}
}

func DeleteATodo(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}

	err := models.DeleteATodo(id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{id: "deleted"})
	}
}
