package controllers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"go.search.crypto/db"
	"go.search.crypto/models"
)

// cleanInput removes every non-alphanumeric character from str.
func sanitizeString(str string) string {
	// Alphanumeric regex.
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		print(err)
		return ""
	}
	return strings.ToLower(reg.ReplaceAllString(str, ""))
}

// Get API endpoint /currency.
// Term query parameter needs to be at least 2 chars long.
func Get(context *gin.Context) {
	queryMap := context.Request.URL.Query() // map[string][]string
	term := queryMap.Get("term")
	searchTerm := sanitizeString(term)
	var currencyList []models.Currency
	if len(searchTerm) < 2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "'term' query param needs to be at least 2 characters long."})
		return
	}
	// Rows where cols: 'name' and 'id' contain the provided term.
	result := db.DB.Where("instr(name, ?) > 0 OR instr(ID, ?) > 0", searchTerm, searchTerm).Find(&currencyList)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": currencyList})
	return
}

// GetAll API endpoint /currency/all
func GetAll(context *gin.Context) {
	var currencyList []models.Currency
	result := db.DB.Find(&currencyList)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": currencyList})
	return
}

// Post API endpoint /currency.
func Post(context *gin.Context) {
	var currency models.Currency
	// Binds and validates request body. Error states which fields are invalid, if any.
	err := context.ShouldBindJSON(&currency)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// New Currency.
	result := db.DB.Create(&currency)
	if result.Error != nil {
		// Returns tx error, if any.
		context.JSON(http.StatusConflict, gin.H{"data": currency, "error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": currency})
	return
}
