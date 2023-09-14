package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	. "example/astar-backend/internal"
)

type RequestJson struct {
	Graph       []Vertex
	Start       Node
	End         Node
	Reruns      int
	Parallelise bool
}

func calculateOptimalPath(c *gin.Context) {
	var requestJson RequestJson
	err := c.BindJSON(&requestJson)
	if err != nil {
		fmt.Println("error")
	}
	var passes int32 = 0 // print it later to make sure all goroutines completed
	var wg sync.WaitGroup
	// complete first pass standalone
	path := CalculatePath(requestJson.Graph, requestJson.Start, requestJson.End, &passes)
	wg.Add(requestJson.Reruns - 1)
	// i starts from 1 because first pass is already completed above
	for i := 1; i < requestJson.Reruns; i++ {
		if requestJson.Parallelise {
			go func() {
				defer wg.Done()
				CalculatePath(requestJson.Graph, requestJson.Start, requestJson.End, &passes)
			}()
		} else {
			wg.Done()
			CalculatePath(requestJson.Graph, requestJson.Start, requestJson.End, &passes)
		}
	}
	wg.Wait()
	c.IndentedJSON(http.StatusOK, path)
}

func homepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("../index.html")
	router.GET("/", homepage)
	router.POST("/astar", calculateOptimalPath)
	router.Run("localhost:8080")
}
