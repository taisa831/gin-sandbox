package main

import (
  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "net/http"
  "strconv"
)

type Todo struct {
  gorm.Model
  Text   string
  Status uint64
}

func main() {
  db, err := gorm.Open("mysql", "gorm:gorm@/sandbox_gin?charset=utf8mb4&parseTime=True&loc=Local")
  if err != nil {
    panic("データベースへの接続に失敗しました")
  }
  defer db.Close()
  db.LogMode(true)

  db.AutoMigrate(&Todo{})

  r := gin.Default()
  r.LoadHTMLGlob("templates/*")

  // 一覧画面
  r.GET("/todo", func(c *gin.Context) {
    var todos []Todo
    db.Find(&todos)

    c.HTML(http.StatusOK, "index.html", gin.H{
      "todos": todos,
    })
  })

  // 新規作成
  r.POST("/todo", func(c *gin.Context) {
    text, _ := c.GetPostForm("text")
    status, _ := c.GetPostForm("status")
    istatus, _ := strconv.ParseUint(status, 10, 32)

    db.Create(&Todo{Text: text, Status: istatus})
    c.Redirect(http.StatusMovedPermanently, "/todo")
  })

  // 編集画面
  r.GET("/todo/:id", func(c *gin.Context) {
    todo := Todo{}
    id := c.Param("id")
    db.First(&todo, id)
    c.HTML(http.StatusOK, "edit.html", gin.H{
      "todo": todo,
    })
  })

  // 編集
  r.POST("/todo/edit/:id", func(c *gin.Context) {
    todo := Todo{}
    id := c.Param("id")
    text, _ := c.GetPostForm("text")
    status, _ := c.GetPostForm("status")
    istatus, _ := strconv.ParseUint(status, 10, 32)

    db.First(&todo, id)
    todo.Text = text
    todo.Status = istatus
    db.Save(&todo)
    c.Redirect(http.StatusMovedPermanently, "/todo")
  })

  // 削除
  r.POST("/todo/delete/:id", func(c *gin.Context) {
    todo := Todo{}
    id := c.Param("id")
    db.First(&todo, id)
    db.Delete(&todo)
    c.Redirect(http.StatusMovedPermanently, "/todo")
  })

  r.Run()
}
