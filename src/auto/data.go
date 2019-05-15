package auto

import "api/models"

var users = []models.User{
	models.User{Nickname: "Jhon Doe", Email: "jhondoe@email.com", Password: "123456789"},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title",
		Content: "Hello world",
	},
}
