/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	pb "github.com/hirasawayuki/go-grpc-cli/pkg/gopher"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	port    = ":8080"
	API_URL = "https://api.github.com"
)

type Server struct {
	pb.UnimplementedGopherServer
}

type Gopher struct {
	URL string `json: "url"`
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the Scheme gRPC server",
	Run: func(cmd *cobra.Command, args []string) {

		listen, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterGopherServer(grpcServer, &Server{})
		log.Printf("gPRC server listening on %v", listen.Addr())

		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func (s *Server) GetGopher(ctx context.Context, req *pb.GopherRequest) (*pb.GopherReply, error) {
	res := &pb.GopherReply{}

	if req == nil {
		fmt.Println("request must not be nil.")
		return res, fmt.Errorf("request must not be nil.")
	}

	if req.Name == "" {
		fmt.Println("name must not be empty in the request")
		return res, fmt.Errorf("name must not be empty in the request.")
	}

	log.Printf("Receive: %v", req.GetName())
	response, err := http.Get(API_URL + "/users/" + req.GetName())

	if err != nil {
		log.Fatalf("failed to call API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("failed to read response body: %v", err)
		}

		var data Gopher
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatalf("failed to unmarshal JSON: %v", err)
		}

		var gophers strings.Builder
		gophers.WriteString(data.URL + "\n")

		res.Message = gophers.String()
	} else {
		log.Fatal("Can't get the Gopher :-(")
	}

	return res, nil
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
