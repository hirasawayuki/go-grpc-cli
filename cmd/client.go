/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/hirasawayuki/go-grpc-cli/pkg/github"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Query the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(address, grpc.WithInsecure())

		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewGithubClient(conn)

		// Github login ID.
		var login string

		if len(os.Args) > 2 {
			login = os.Args[2]
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := client.GetGithubUser(ctx, &pb.GithubUserRequest{Login: login})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Github user URL: %s", r.GetHtmlUrl())
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
