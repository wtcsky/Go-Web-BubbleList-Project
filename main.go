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
