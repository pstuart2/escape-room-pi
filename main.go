package main

import (
	"log"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"

	"github.com/labstack/echo"
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"time"
	"context"
)

const (
	Light1 int = 2
	Light2     = 3
	Light3     = 4
	Light4     = 17
	Light5     = 27
	Light6     = 22
	Light7     = 10
	Light8     = 9
)

var pins []embd.DigitalPin

func main() {
	log.Println("running...")

	lights := []int{Light1, Light2, Light3, Light4, Light5, Light6, Light7, Light8}

	embd.InitGPIO()
	defer embd.CloseGPIO()

	for _, v := range lights {
		pin, err := embd.NewDigitalPin(v)
		if err != nil {
			log.Fatal(err)
		}

		pin.SetDirection(embd.Out)
		pin.Write(embd.High)
		pins = append(pins, pin)
	}

	e := echo.New()
	e.POST("/lights", Lights)

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func Lights(ctx echo.Context) error {
	r := new([]uint)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	fmt.Printf("LightRequest: %v\n", r)
	for i, v := range *r {
		if v == 0 {
			setLightOff(i)
		} else {
			setLightOn(i)
		}
	}

	return nil
}

func setLightOn(i int) {
	pins[i].Write(embd.Low)
}

func setLightOff(i int) {
	pins[i].Write(embd.High)
}
