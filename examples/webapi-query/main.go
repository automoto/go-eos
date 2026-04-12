// Command webapi-query demonstrates the pure-Go EOS Web API client.
//
// This program builds with CGO_ENABLED=0 and requires no native EOS SDK.
// It authenticates via client_credentials and queries leaderboard definitions.
//
// Environment variables:
//
//	EOS_CLIENT_ID       — from the Epic Developer Portal
//	EOS_CLIENT_SECRET   — from the Epic Developer Portal
//	EOS_DEPLOYMENT_ID   — from the Epic Developer Portal
//
// Usage:
//
//	CGO_ENABLED=0 go run ./examples/webapi-query
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mydev/go-eos/webapi"
)

func main() {
	clientID := requireEnv("EOS_CLIENT_ID")
	clientSecret := requireEnv("EOS_CLIENT_SECRET")
	deploymentID := requireEnv("EOS_DEPLOYMENT_ID")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := webapi.New(deploymentID,
		webapi.WithClientCredentials(clientID, clientSecret),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("Authenticated via client_credentials.")

	// Query leaderboard definitions.
	fmt.Println("\n--- Leaderboard Definitions ---")
	defs, err := client.GetLeaderboardDefinitions(ctx)
	if err != nil {
		log.Fatalf("GetLeaderboardDefinitions: %v", err)
	}
	if len(defs) == 0 {
		fmt.Println("(no leaderboards configured for this deployment)")
	}
	for _, d := range defs {
		fmt.Printf("  %s (stat: %s, aggregation: %s)\n",
			d.Spec.Name, d.Spec.RankBy.Stat, d.Spec.RankBy.Aggregation)
		if len(d.Players) == 0 {
			fmt.Println("    (no players)")
		}
		for _, e := range d.Players {
			fmt.Printf("    #%d  %s  score=%d\n", e.Rank, e.ProductUserID, e.Score)
		}
	}
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return val
}
