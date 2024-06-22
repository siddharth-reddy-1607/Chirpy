package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type AppMetrics struct{
    TimesVisited int
}

func (am *AppMetrics) middlwareUpdateMetrics(next http.Handler) http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        am.TimesVisited += 1 
        next.ServeHTTP(w,r)
    })
}

func (am *AppMetrics) displayMetrics() http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        w.Header().Add("Content","text/html")
        templ,err := template.ParseFiles("./admin/metrics.html")
        if err!= nil{
            err_msg := fmt.Sprintf("Error while parsing files for metrics.html : %s",err)
            http.Error(w, err_msg, http.StatusInternalServerError)
            return
        }
        if err := templ.Execute(w,am); err!= nil{
            err_msg := fmt.Sprintf("Error while executine the template for metrics.html : %s",err)
            http.Error(w, err_msg, http.StatusInternalServerError)
            return
        }
    })
}

func (am *AppMetrics) resetMetrics() http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        am.TimesVisited = 0
        w.WriteHeader(http.StatusOK)
    })
}
