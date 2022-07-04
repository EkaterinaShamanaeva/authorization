package main

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/autorization/internal/application"
	"github.com/EkaterinaShamanaeva/autorization/internal/repository"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	// подключение к БД
	dbpool, err := repository.InitDbConn(ctx)
	if err != nil {
		log.Fatalf("%w failed to init DB connection", err)
	}
	defer dbpool.Close()

	a := application.NewApp(ctx, dbpool)
	// роутер
	r := httprouter.New()
	a.Routes(r)
	srv := &http.Server{Addr: "0.0.0.0:8070", Handler: r}
	fmt.Println("It is alive! Try http://localhost:8070")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("%w failed to connect", err)
	}
}
