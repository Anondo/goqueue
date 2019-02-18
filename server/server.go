package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

func Serve() error {
	port := viper.GetInt("port")

	prepareRouter()

	r := chi.NewMux()
	r.Mount("/api/", router)

	return http.ListenAndServe(":"+strconv.Itoa(port), r)

}
