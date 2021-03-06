package controller

import (
	"log"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwiksqlite-admin/helpers"
	"github.com/jasonbronson/kwiksqlite-admin/repository"
	"gorm.io/gorm"
)

func ConnectDB(g *gin.Context) {
	//helpers.Cfg.DbName = g.GetHeader("database")
	helpers.Cfg.DbName = g.Query("db")
	if !repository.CheckDBExists(helpers.Cfg.DbName) {
		g.JSON(500, gin.H{"error": "database file does not exist"})
		return
	}
	d := repository.GetDatabaseInfo()
	g.JSON(200, d)
}

func GetDatabaseInfo(g *gin.Context) {
	//TODO fix issue with missing dbstats
	d := repository.GetDatabaseInfo()
	g.JSON(200, d)
}

func CustomQuery(g *gin.Context) {

	json := struct {
		Query string `json:"query"`
	}{
		"",
	}
	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(500, gin.H{"error": "query is required"})
		return
	}

	log.Println(json.Query)

	var result interface{}
	err := helpers.DB().Raw(json.Query).Scan(&result).Error
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, result)
}

func GetTables(g *gin.Context) {
	var result []ShowTables
	helpers.DB().Raw("SELECT name, sql FROM sqlite_master WHERE type ='table' AND name NOT LIKE 'sqlite_%' order by name").Scan(&result)
	g.JSON(200, result)
}

func DropTable(g *gin.Context) {

	table := g.Param("tablename")
	//check if table exists
	if table == "" {
		g.JSON(500, gin.H{"error": "table name is required"})
		return
	}
	err := helpers.DB().Debug().Migrator().DropTable(table)
	if err != nil {
		log.Println(err)
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	d := repository.GetDatabaseInfo()
	g.JSON(200, d)
}

func GetTableContent(g *gin.Context) {
	table := g.Param("tablename")
	c, err := repository.GetTableContent(table)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Println(c)
	g.JSON(200, c)
}

func GetColumns(g *gin.Context) {
	table := g.Param("tablename")
	c := repository.GetColumns(table)
	g.JSON(200, c)
}

func CreateTable(g *gin.Context) {
	table := g.Param("tablename")
	match, _ := regexp.MatchString("[a-zA-Z_]([0-9]?)", table)
	if !match {
		g.JSON(500, gin.H{"error": "table name does not match table name conventions"})
		return
	}
	helpers.DB().Table(table).AutoMigrate(&NewTable{})
	d := repository.GetDatabaseInfo()
	g.JSON(200, d)
}

type ShowTables struct {
	Name string
	SQL  string
}

type NewTable struct {
	gorm.Model
}
