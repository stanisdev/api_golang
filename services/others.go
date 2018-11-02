package services

import (
	"github.com/spf13/viper"
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