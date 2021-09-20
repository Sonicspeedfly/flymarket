package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Sonicspeedfly/flymarket/cmd/app"
	"github.com/Sonicspeedfly/flymarket/pkg/accounts"
	"github.com/Sonicspeedfly/flymarket/pkg/product"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/db"

	if err := execute(host, port, dsn); err != nil {
		os.Exit(1)
	}
}

func HTTPServer(server *app.Server, host, port string) *http.Server {
	return &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
}

func execute(host string, port string, dsn string) (err error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
	defer func(){
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()
	mux := mux.NewRouter()
	accountsSvc := accounts.NewService(db)
	productSvc := product.NewService(db)
	server := app.NewServer(mux, accountsSvc, productSvc)
	log.Println("server start")

	a := func(server *app.Server) *http.Server {
		server.Init()
		return &http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: server,
		}
	}(server)
	if err != nil {
		log.Fatal("http server: ", err)
	}
	a.ListenAndServe()
	return err
}

//log.Println("handler{}.ServeHTTP()")
//writer.Write([]byte("hello http"))
