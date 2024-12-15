package main

import (
	"net/http"

	"backend/algorithms"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// API 路由
	api := r.Group("/api")
	{
		api.POST("/prim", handlePrim)
		api.POST("/dijkstra", handleDijkstra)
		api.POST("/floyd", handleFloyd)
		api.POST("/bayes", handleBayes)
		api.POST("/max-matching", handleMaxMatching)
		api.POST("/tsp", handleTSP)
	}

	r.Run(":8080")
}

// 处理 Prim 算法请求
func handlePrim(c *gin.Context) {
	var input struct {
		Graph [][]int `json:"graph"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := algorithms.Prim(input.Graph)
	c.JSON(http.StatusOK, gin.H{"result": result})
}

// 处理 Dijkstra 算法请求
func handleDijkstra(c *gin.Context) {
	var input struct {
		Graph [][]int `json:"graph"`
		Start int     `json:"start"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := algorithms.Dijkstra(input.Graph, input.Start)
	c.JSON(http.StatusOK, gin.H{"result": result})
}

// 处理 Floyd 算法请求
func handleFloyd(c *gin.Context) {
	var input struct {
		Graph [][]int `json:"graph"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := algorithms.Floyd(input.Graph)
	c.JSON(http.StatusOK, gin.H{"result": result})
}

// 处理 Bayes 算法请求
func handleBayes(c *gin.Context) {
	var input struct {
		Network  algorithms.BayesNetwork `json:"network"`
		Evidence map[int]bool            `json:"evidence"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := algorithms.Bayes(input.Network, input.Evidence)
	c.JSON(http.StatusOK, gin.H{"result": result})
}

// 处理最大匹配算法请求
func handleMaxMatching(c *gin.Context) {
	var input struct {
		Graph [][]int `json:"graph"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := algorithms.MaxMatching(input.Graph)
	c.JSON(http.StatusOK, gin.H{"result": result})
}

// 处理旅行商问题算法请求
func handleTSP(c *gin.Context) {
	var input struct {
		Graph [][]int `json:"graph"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := algorithms.TSP(input.Graph)
	c.JSON(http.StatusOK, gin.H{"result": result})
}
