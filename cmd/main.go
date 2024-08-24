package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
}

func main() {

	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)

	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.Authenticator{}.Client(ctx, token)
	client := spotify.New(httpClient)

	client.
		results, err := client.Search(ctx, "christmas", spotify.SearchTypePlaylist)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(results.Playlists)
}
