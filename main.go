package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/joaquinicolas/xmen"
	"fmt"
	"os"
)

type XMen struct {
	DNA []string `json:"dna" binding:"required"`
}
func main() {
	router := gin.Default()
	router.POST("/mutant", func(c *gin.Context) {
		var mutant XMen
		if err := c.BindJSON(&mutant); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err})
			return
		}

		rowsLength := len(mutant.DNA)
		colsLength := len(mutant.DNA[0])
		var dna [][]string
		for x := 0; x < rowsLength; x++ {
			var cols []string
			for i := 0; i < colsLength; i++ {
				cols = append(cols, string(mutant.DNA[x][i]))
			}
			fmt.Println(cols)
			dna = append(dna, cols)
		}

		isMutant := xmen.IsMutant(dna)
		c.JSON(http.StatusOK, gin.H{"isMutant": isMutant})
	})

	port := os.Getenv("PORT") // Heroku provides the port to bind to
	if port == "" {
		port = "8080"
	}

	router.Run(os.Getenv(":"+port))
}
