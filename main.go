package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"github.com/joaquinicolas/xmen"
	"fmt"
	"strings"
	"crypto/sha1"
)

type XMen struct {
	DNA []string `json:"dna" binding:"required"`
}

func main() {
	db, err := sql.Open("sqlite3", "./xmen.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	createTable(db)

	router := gin.Default()
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK,"It's works!")
	})
	router.GET("/stats", getStats(db))
	router.POST("/mutant", postMutant(db))

	router.Run()
}

func postMutant(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			dna = append(dna, cols)
		}

		isMutant := xmen.IsMutant(dna)
		err := saveDNA(mutant.DNA, isMutant, db)
		if err != nil {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"isMutant": isMutant})
	}
}

func getStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var mutantsCounter float32
		err := db.QueryRow("SELECT COUNT(ID) FROM DNAS WHERE isMutant;").Scan(&mutantsCounter)
		if err != nil {
			c.String(http.StatusForbidden, err.Error())
			return
		}

		var humansCounter float32
		err = db.QueryRow("SELECT COUNT(ID) FROM DNAS WHERE NOT isMutant;").Scan(&humansCounter)
		if err != nil {
			c.String(http.StatusForbidden, err.Error())
			return
		}

		var ratio float32 = 0
		if humansCounter > 0 {
			ratio = mutantsCounter / humansCounter
		}

		c.JSON(http.StatusOK, gin.H{"count_mutant_dna": mutantsCounter, "count_human_dna":humansCounter, "ratio": ratio})
	}
}

func saveDNA(dna []string, isMutant bool, db *sql.DB) error {
	dnaJoined := strings.Join(dna, ",")
	hash := sha1.New()
	hash.Write([]byte(dnaJoined))
	bs := hash.Sum(nil)
	query := "INSERT INTO DNAS(ID,dna, isMutant) VALUES (?, ?, ?)"

	stmt, _ := db.Prepare(query)
	_, err := stmt.Exec(fmt.Sprintf("%x", bs), dnaJoined, isMutant)

	return err
}


func createTable(db *sql.DB) error {
	query := `CREATE TABLE dnas
(
    ID TEXT PRIMARY KEY,
    dna TEXT,
	isMutant BOOLEAN DEFAULT FALSE  NOT NULL
);
CREATE UNIQUE INDEX dnas_ID_uindex ON dnas (ID);`

	_, err := db.Exec(query)
	return err
}