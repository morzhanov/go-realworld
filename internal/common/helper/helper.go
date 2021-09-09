package helper

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
)

func CheckStruct(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Struct
}

func HandleInitializationError(err error, step string, log *zap.Logger) {
	log.Fatal(
		fmt.Sprintf("an error occurred during %s initialization step: %s", step, err),
		zap.String("step", step),
	)
}
