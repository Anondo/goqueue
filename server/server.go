package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

// Serve serves the http server for goqueue
func Serve() error {
	port := viper.GetInt("port")

	prepareRouter()

	r := chi.NewMux()
	r.Mount("/api/", router)

	return http.ListenAndServe(":"+strconv.Itoa(port), r)

}
