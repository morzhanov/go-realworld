package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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

func StartGrpcServer(s *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	log.Printf("Grpc server started at: localhost%v", port)
	return nil
}
