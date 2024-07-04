package main

import (
    // "fmt"
    // "log"
    "net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request){
	respondWithJSON(w, 200, struct{}{} )
}