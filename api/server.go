package main

import (
	"hello/server/infra"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 出力の定義
// @title Techer API
// @version 1.0
// @description This is a Api server
// @license.name Apache 2.0
// @BasePath /
func main() {
	e := echo.New() //eにecho.New(echoを使っている)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost", "http://quiztecher.com"},
		AllowMethods: []string{
			http.MethodGet, http.MethodHead,
			http.MethodPut, http.MethodOptions,
			http.MethodPatch, http.MethodPost,
			http.MethodDelete},
		AllowHeaders: []string{
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderContentType,
			echo.HeaderContentLength,
			echo.HeaderAcceptEncoding,
			echo.HeaderXCSRFToken,
			echo.HeaderAuthorization},
		AllowCredentials: true,
	})) //echoはCORSを使う
	infra.Routing(e)                 //infraのRoutingを実行
	e.Logger.Fatal(e.Start(":8080")) //ポート解放している
}
