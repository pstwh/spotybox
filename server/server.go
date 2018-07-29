package main

import (
	"context"
	"fmt"
	"os"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	config = &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, _ = config.Token(context.Background())
	client = spotify.Authenticator{}.NewClient(token)
	upgrader = websocket.Upgrader{}
)

func hello(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for{
		header, page, _ := client.FeaturedPlaylists()

		var msg string
		msg = "<hr>" + msg + header + "<br>"

		for _, playlist := range page.Playlists {
			msg = msg + playlist.Name + "<br>"
		}

		err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			c.Logger().Error(err)
		}

		_, m, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", m)
	}
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "../public")
	e.GET("/ws", hello)
	e.Logger.Fatal(e.Start(":1323"))
}
