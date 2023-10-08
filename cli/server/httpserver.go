package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"serverless/internal/engine"
	"serverless/internal/executor"
	"time"
)

type ServerlessHttpServer interface {
	Start() error
	Stop(ctx context.Context) error
}

type InternalHttpServer struct {
	listen     string
	port       int
	folder     string
	manifest   *engine.Manifest
	httpServer *http.Server
}

func NewInternalHttpServer(listen string, port int, folder string, manifest *engine.Manifest) ServerlessHttpServer {
	return &InternalHttpServer{listen: listen, port: port, folder: folder, manifest: manifest}
}

func (s *InternalHttpServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	responseChannel, err := executor.Execute(s.folder, s.manifest)
	if err != nil {
		log.Print(err)
	}
	response := <-responseChannel

	// Write Headers
	for key, value := range response.Headers {
		print(key + ":" + fmt.Sprintf("%v", value) + "\n")
		rw.Header().Set(key, fmt.Sprintf("%v", value))
	}
	rw.WriteHeader(response.Status)

	// Write Response
	_, err = rw.Write([]byte(response.Content))
	if err != nil {
		return
	}
}

func (s *InternalHttpServer) Start() error {
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.listen, s.port),
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.httpServer.ListenAndServe()
}

func (s *InternalHttpServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
