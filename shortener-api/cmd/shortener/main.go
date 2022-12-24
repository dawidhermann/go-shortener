package main

import (
	"github.com/dawidhermann/shortener-api/api"
	"github.com/dawidhermann/shortener-api/internal/db"
	"os"
)

//import pb "github.com/dawidhermann/shortener-api/internal/protobuf"

func main() {
	//url := shortener.ShortenUrl()
	//fmt.Println(url)
	//api.StartServer()
	//db.Connect("shortener_user", "P0sTgr3sP4SS", "postgres:5432", "shortener_db")
	db.Connect(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), "postgres:5432", os.Getenv("POSTGRES_DB"))
	//db.CreateUser("dhermann", "secretpass", "email@example.com", dbInstance)

	api.StartServer()

}
