# Go-Web-BubbleList-Project
通过Gin框架以及GORM来编写的一个小清单的项目
## 0. 描述

一个前后端分离的小清单项目，我只负责后端功能的实现，前端的一些设计(包括js、css、html等)直接从[Github](https://github.com/Q1mi/go_web/tree/master/lesson25)上下载即可。

通过使用Gin框架来进行与客户端(网页)的信息交互以及GORM来实现服务端与数据库的数据交互(CRUD)。

需要完成的功能是显示数据库中的所有待办事项、添加一条新的待办事项、将未完成的待办事项设定为已完成、将已完成的待办事项重置为未完成、删除指定的待办事项。

完整代码下载

## 1. 项目展示

### 1.1 Web页面效果展示

![chrome_afanK9rv5z.gif](https://s2.loli.net/2023/03/15/QiDaAnEIXtFMBvW.gif)

### 1.2 MySQL效果展示

![navicat_3mYCFponOo.gif](https://s2.loli.net/2023/03/15/aYMvqbxl2sXK3LC.gif)

## 2. 项目介绍

项目的基本结构如下图所示：

<img src="https://s2.loli.net/2023/03/15/DaUiGJg9VbPXdwY.png" alt="goland64_KDAreF6nbQ.png" style="zoom:50%;" />

### 2.1 持久层DAO

这一部分就是通过使用`grom`包内的方法与MySQL数据库进行连接，具体如下：

```go
package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitMySQL() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bubble_list?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}
```

这里的dsn设置中，`root`是用户名，`123456`是密码，在`127.0.0.1`(本机)的`3306`端口建立`tcp`连接，使用的数据库名是`bubble_list`，之后是其他的一些连接配置。

在调用了`dao.InitMySQL()`之后，对于数据库的操作都可以通过`dao.DB`直接使用。

### 2.2 模型Model

### 2.2.1 Todo模型

因为要完成的是一个清单列表，因此要操作的对象就应该是一条待办事项(包括唯一标识ID、待办事项的具体内容、是否完成)，很容易就能想到ID可以用一个`int`类型；具体内容就是`string`类型；至于是否完成了该条事项就用一个`bool`类型就可以了。

于是我们就顺理成章的得到了下面的结构体：

```go
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}
```

可以看到其实这里还在每一条后面带了一个`tag`，这是为了与前端做数据交流做准备的。

成员变量的首字母不大写其他的包无法访问，可能会导致一些其他包的方法无法正常使用，那之后通过使用`Gin`框架的方法与请求中的`JSON`进行数据交互的时候结构体中的成员变量需要与`JSON`中的信息对应，因此需要加上这些`tag`。

### 2.1.2 Todo模型的CRUD操作

```go
func CreateATodo(todo *Todo) (err error) {
	err = dao.DB.Create(&todo).Error
	return
}

func GetTodoList(todoList *[]Todo) (err error) {
	err = dao.DB.Find(&todoList).Error
	return
}

func GetATodo(id string) (todo *Todo, err error) {
	if err = dao.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		return nil, err
	}
	return
}
func UpdateATodo(todo *Todo) (err error) {
	err = dao.DB.Save(&todo).Error
	return
}

func DeleteATodo(id string) (err error) {
	err = dao.DB.Where("id = ?", id).Delete(&Todo{}).Error
	return
}
```

以上操作都是调用了`grom`包中的方法来进行的。

### 2.3 控制器Controller

由于执行的操作逻辑比较简单因此省略了业务逻辑层Logic。这里是用来存放处理客户端请求的具体操作的。其中对于Todo模型的CRUD是直接通过调用`models`包里之前已经写好的内容来完成的。

```go
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
```

主要也是一些Gin框架的内容，

### 2.4 路由Router

这里就是通过使用`Gin`框架来完成一些对页面请求的一些处理，具体的处理方法放在了前面的`Controller`层。以下是创建了一个gin的引擎之后通过该引擎进行操作设置静态文件位置、加载模板文件、定义路由以及定义路由组。

```go
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
```

通过不同的请求方式完成不同的需求，比如`POST`请求是增加新的数据，`GET`请求是查询数据，`PUT`请求是更新数据，`Delete`请求是删除数据……

### 2.5 主函数main

```go
package main

import (
   "BubbleList/dao"
   "BubbleList/models"
   "BubbleList/routers"
)

// 连接数据库

func main() {
   // 连接数据库
   err := dao.InitMySQL()
   if err != nil {
      panic(err)
   }
   // 绑定Model
   _ = dao.DB.AutoMigrate(&models.Todo{})
   r := routers.SetupRouter()
   _ = r.Run()
}
```

这里就是一个宏观上的使用了，也是最接近客户端的位置。

### 3. 参考链接

[1]: https://www.bilibili.com/video/BV1gJ411p7xC	"基于gin框架和gorm的web开发实战 (七米出品)"

