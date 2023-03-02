package main

import (
	"WpeWebui/src/wpe"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	workshopPath := "D:\\SteamLibrary\\steamapps\\workshop\\content\\431960"
	index := wpe.New(workshopPath)

	app := gin.New()
	app.LoadHTMLGlob("src/html/**")
	app.GET("/index", func(context *gin.Context) {
		context.HTML(200, "index.html", gin.H{
			"projects": index.List(),
		})
	})
	app.GET("/project/:id", func(context *gin.Context) {
		id := context.Param("id")
		fullPath := index.FullPath(id)
		open, err := os.Open(fullPath)
		if err != nil {
			context.String(404, "id: %s, not fond", id)
		}
		stat, err := open.Stat()
		if err != nil {
			context.String(404, "id: %s, not fond", id)
		}
		http.ServeContent(context.Writer, context.Request, "test.mp4", stat.ModTime(), open)
	})
	port := 8080
	if err := app.Run(fmt.Sprintf(":%d", port)); err != nil {
		return
	}
}
