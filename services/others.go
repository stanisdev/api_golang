package services

import (
	"github.com/spf13/viper"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LimitRestriction(limit int, entityName string) int {
	maxLimit, err := strconv.Atoi(viper.GetString("limiter_per_page." + entityName))
	if (err != nil) {
		panic("The entityName " + entityName + " in config file, related with limiter section, does not exist")
	}
	if (limit > maxLimit) {
		limit = maxLimit
	}
	return limit
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
			header["Content-Type"] = value
	}
}

func JSONgoesToHTML(c *gin.Context, obj interface{}) {
	c.Status(200)
	writeContentType(c.Writer, []string{"application/json; charset=utf-8"})
	enc := json.NewEncoder(c.Writer)
	enc.SetEscapeHTML(false) // this is missing in gin
	if err := enc.Encode(obj); err != nil {
			panic(err)
	}
}