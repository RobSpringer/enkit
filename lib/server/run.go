package server

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(mux *http.ServeMux, grpcs *grpc.Server) {
	reflection.Register(grpcs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6433"
	}

	log.Printf("Opening port %s - will be available at http://127.0.0.1:%s/", port, port)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	// Create all listeners.
	cml := cmux.New(listener)
	grpcl := cml.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpl := cml.Match(cmux.Any())

	grpcw := grpcweb.WrapServer(grpcs, grpcweb.WithAllowNonRootResource(true), grpcweb.WithWebsockets(true), grpcweb.WithOriginFunc(func(string) bool { return true }))

	https := &http.Server{Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if grpcw.IsGrpcWebRequest(req) {
			grpcw.ServeHTTP(resp, req)
		} else {
			mux.ServeHTTP(resp, req)
		}
	})}
	go grpcs.Serve(grpcl)
	go https.Serve(httpl)

	if err := cml.Serve(); err != nil {
		log.Fatalf("Serve failed with error: %s", err)
	}
}
