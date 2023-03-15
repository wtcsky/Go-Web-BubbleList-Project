package models

import "BubbleList/dao"

// Todo Model
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

// Todo这个Model的CRUD操作都放在这里

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
