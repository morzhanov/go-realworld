package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
func StartGrpcServer(ctx context.Context, s *grpc.Server, port string, log *zap.Logger) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	log.Info("Grpc server started", zap.String("port", port))
	<-ctx.Done()
	lis.Close()
	return nil
}

func HandleInitializationError(err error, step string, log *zap.Logger) {
	log.Fatal(
		fmt.Sprintf("an error occurred during %step initialization step: %s", step, err),
		zap.String("step", step),
	)
}

// TODO: maybe we should move this somwhere else like in BaseEventsController
func StartRestServer(ctx context.Context, port string, router *gin.Engine, log *zap.Logger) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("REST Server Failed to start", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutdown REST Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("REST Server shutdown failed", zap.Error(err))
	}
}
