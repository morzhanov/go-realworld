package helper

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
)

func CheckStruct(val interface{}) bool {
	kind := reflect.ValueOf(val).Kind()
	return kind == reflect.Struct || kind == reflect.Ptr
}

func HandleInitializationError(err error, step string, log *zap.Logger) {
	log.Fatal(
		fmt.Sprintf("an error occurred during %s initialization step: %s", step, err),
		zap.String("step", step),
	)
}

func CheckNotFound(err error) bool {
	return err.Error() == "sql: no rows in result set"
}
