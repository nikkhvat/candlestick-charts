package server

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"forex/models"
	"forex/pkg/config"
	database "forex/pkg/database"
	"forex/pkg/middleware"

	gin "github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"forex/terminal"
)

type App struct {
	httpServer *http.Server
	db         *gorm.DB
}

func NewApp() *App {
	db := database.InitDB()

	return &App{
		db: db,
	}
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients []*websocket.Conn

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	clients = append(clients, conn)

	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		for _, client := range clients {
			client.WriteMessage(t, msg)
		}
	}
}

func (a *App) Run(port string) error {
	conf := config.GetConfig()

	flag.Parse()
	log.SetFlags(0)

	messageOut := make(chan string)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: "marketdata.tradermade.com", Path: "/feedadv"}

	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var data models.Data

			json.Unmarshal([]byte(string(message)), &data)

			for _, client := range clients {
				client.WriteMessage(1, message)
			}

			a.db.Create(data)

			if string(message) == "Connected" {
				log.Printf("Send Sub Details: %s", message)

				key := conf.Key
				nominals := conf.Nominals

				message := "{\"userKey\":\"" + key + "\", \"symbol\":\"" + nominals + "\"}"

				messageOut <- message
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("/forex/ping", func(c *gin.Context) {
		c.String(200, "ping")
	})

	router.GET("/forex/actual", func(c *gin.Context) {
		symbol := c.Query("symbol")

		c.JSON(200, terminal.ActualDate(a.db, symbol))
	})

	router.GET("/forex/hostory", func(c *gin.Context) {
		start_date := c.Query("start_date")
		end_date := c.Query("end_date")
		symbol := c.Query("symbol")
		timeframe := c.Query("timeframe")

		start_date_int64, err := strconv.Atoi(start_date)

		if err != nil {
			c.JSON(400, "err 1")
			return
		}

		end_date_int64, err := strconv.Atoi(end_date)

		if err != nil {
			c.JSON(400, "err 2")
			return
		}

		timeframe_int, err := strconv.Atoi(timeframe)

		if err != nil {
			c.JSON(400, "err 3")
			return
		}

		candles := terminal.Hostory(a.db, symbol, int64(start_date_int64), int64(end_date_int64), timeframe_int)

		c.JSON(http.StatusOK, candles)
	})

	router.GET("/forexws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case m := <-messageOut:
				log.Printf("Send Message %s", m)
				err := c.WriteMessage(websocket.TextMessage, []byte(m))
				if err != nil {
					log.Println("write:", err)
					return
				}
			case t := <-ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
				if err != nil {
					log.Println("write:", err)
					return
				}
			case <-interrupt:
				log.Println("interrupt")

				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
