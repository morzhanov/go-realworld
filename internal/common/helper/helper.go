package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"time"

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

// TODO: maybe we should move this somwhere else like in BaseEventsController
func StartGrpcServer(ctx context.Context, s *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	log.Printf("Grpc server started at: localhost%v", port)
	<-ctx.Done()
	lis.Close()
	return nil
}

func HandleInitializationError(err error, step string) {
	log.Fatal(
		fmt.Errorf("an error occurred during %step initialization step: %w", step, err),
	)
}

// TODO: maybe we should move this somwhere else like in BaseEventsController
func StartRestServer(ctx context.Context, port string, router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// TODO: refactor with appropriate error handling
	<-ctx.Done()
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
