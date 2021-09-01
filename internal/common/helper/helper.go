package helper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func CheckStruct(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Struct
}

func ParseRestBody(ctx *gin.Context, input interface{}) error {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	in := reflect.ValueOf(input)
	return json.Unmarshal(jsonData, &in)
}

func HandleRestError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}
