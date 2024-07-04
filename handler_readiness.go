package main

import (
    // "fmt"
    // "log"
    "net/http"
)

func (app *App) handlerReadiness(w http.ResponseWriter, r *http.Request){
	respondWithJSON(w, 200, struct{}{} )
}