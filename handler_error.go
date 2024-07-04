package main

import (
	"net/http"
)

func (app *App) handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Someting Went Wrong!!")
}
