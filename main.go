package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vitaly-kashtalyan/hlk-sw16"
	"net/http"
	"os"
	"strconv"
	"time"
)

type httpError struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

type httpOkData struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data" example:"interface{}"`
}

type Relay struct {
	Id    int `json:"id"`
	State int `json:"state"`
}

func init() {
	if len(os.Getenv("HLK_SW16_HOST")) == 0 {
		_ = os.Setenv("HLK_SW16_HOST", "192.168.0.200")
	}

	if len(os.Getenv("HLK_SW16_PORT")) == 0 {
		_ = os.Setenv("HLK_SW16_PORT", "8080")
	}

	if len(os.Getenv("APP_PORT")) == 0 {
		_ = os.Setenv("APP_PORT", "8082")
	}
}

func main() {
	// Runs the server
	r := setup()
	_ = r.Run(":" + os.Getenv("APP_PORT")) // listen and serve on 0.0.0.0:8082
}

func setup() *gin.Engine {
	fmt.Println("Starting routes with gin")
	r := gin.Default()
	initializeRoutes(r)
	return r
}

func initializeRoutes(r *gin.Engine) {
	r.GET("/", getNewStatus)
	r.GET("/status", getStatus)
	r.GET("/relays/on", switchOnAll)
	r.GET("/relays/on/:id", switchOn)
	r.GET("/relays/off", switchOffAll)
	r.GET("/relays/off/:id", switchOff)
}

func switchOff(c *gin.Context) {
	id, err := getUInt(c.Param("id"))
	if err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}

	hlk := getConnect()
	if hlk.Err != nil {
		resError(c, http.StatusServiceUnavailable, hlk.Err)
		return
	}

	if err := hlk.RelayOff(id); err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}

	msg, err := hlk.ReadMessage()
	if err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}
	_ = hlk.Close()
	resOkData(c, setMapRelays(msg))
}

func switchOn(c *gin.Context) {
	id, err := getUInt(c.Param("id"))
	if err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}

	hlk := getConnect()
	if hlk.Err != nil {
		resError(c, http.StatusServiceUnavailable, hlk.Err)
		return
	}

	if err := hlk.RelayOn(id); err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}

	msg, err := hlk.ReadMessage()
	if err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}
	_ = hlk.Close()
	resOkData(c, setMapRelays(msg))
}

func switchOffAll(c *gin.Context) {
	hlk := getConnect()
	if hlk.Err != nil {
		resError(c, http.StatusServiceUnavailable, hlk.Err)
		return
	}

	if err := hlk.SwitchAllOff(); err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}

	msg, err := hlk.ReadMessage()
	if err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}
	_ = hlk.Close()
	resOkData(c, setMapRelays(msg))
}

func switchOnAll(c *gin.Context) {
	hlk := getConnect()
	if hlk.Err != nil {
		resError(c, http.StatusServiceUnavailable, hlk.Err)
		return
	}

	if err := hlk.SwitchAllOn(); err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}

	msg, err := hlk.ReadMessage()
	if err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}
	_ = hlk.Close()
	resOkData(c, setMapRelays(msg))
}

func getNewStatus(c *gin.Context) {
	for i := 0; i < 10; i++ {
		hlk := getConnect()
		if hlk.Err != nil {
			resError(c, http.StatusServiceUnavailable, hlk.Err)
			return
		}

		if err := hlk.StatusRelays(); err != nil {
			resError(c, http.StatusBadRequest, err)
			return
		}

		msg, err := hlk.ReadMessage()
		if err != nil {
			resError(c, http.StatusBadRequest, err)
			return
		}

		if validateAnswer(msg) {
			resOkData(c, setNewMapRelays(msg))
			return
		} else {
			time.Sleep(1000)
		}
		_ = hlk.Close()
	}
}

func getStatus(c *gin.Context) {
	hlk := getConnect()
	if hlk.Err != nil {
		resError(c, http.StatusServiceUnavailable, hlk.Err)
		return
	}

	if err := hlk.StatusRelays(); err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}

	msg, err := hlk.ReadMessage()
	if err != nil {
		resError(c, http.StatusBadRequest, err)
		return
	}
	_ = hlk.Close()
	resOkData(c, setMapRelays(msg))
}

func getConnect() (c *hlk_sw16.Connection) {
	return hlk_sw16.New(os.Getenv("HLK_SW16_HOST"), os.Getenv("HLK_SW16_PORT"))
}

func setNewMapRelays(msg []byte) (relay []Relay) {
	for index, element := range msg {
		if index > 1 && index < 18 {
			status := int(element)
			if status == 2 {
				status = 0
			}
			relay = append(relay, Relay{
				Id:    int(index) - 2,
				State: status,
			})
		}
	}
	return
}

func setMapRelays(msg []byte) (relays map[int]int) {
	relays = make(map[int]int)
	for index, element := range msg {
		if index > 1 && index < 18 {
			status := int(element)
			if status == 2 {
				status = 0
			}
			relays[int(index)-2] = status
		}
	}
	return
}

func validateAnswer(msg []byte) bool {
	for index, element := range msg {
		if index > 1 && index < 18 {
			if int(element) > 2 {
				return false
			}
		}
	}
	return true
}

func resError(ctx *gin.Context, code int, message error) {
	er := httpError{
		Status:  code,
		Message: message.Error(),
	}
	ctx.JSON(http.StatusOK, er)
}

func resOkData(ctx *gin.Context, data interface{}) {
	ok := httpOkData{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, ok)
}

func getUInt(value string) (int, error) {
	i, err := strconv.ParseInt(value, 10, 32)
	return int(i), err
}
