package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct{
	ID string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL string `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`

}
var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher:=md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("hasher: ",hasher)
	data:= hasher.Sum(nil)
	fmt.Println("data: ",data)
	hash:=hex.EncodeToString(data)
	fmt.Println("hash: ",hash)
	fmt.Println("final string: ",hash[:8])
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL:=generateShortURL(originalURL)
	id:=shortURL
	urlDB[id]=URL{
		ID: id,
		OriginalURL: originalURL,
		ShortURL: shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}
func rootURLHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("GET function")
	fmt.Fprintf(w,"hello World")
}

func shortURLHandler(w http.ResponseWriter, r *http.Request){
	var data struct{
		URL string `json:"url"`
	}
	err:=json.NewDecoder(r.Body).Decode(&data)
	if err!=nil{
		http.Error(w,"Invalid Request body", http.StatusBadRequest)
		return
	}

	shortURL_:=createURL(data.URL)
	//fmt.Fprintf(w,shortURL)
	response:=struct{
		ShortURL string `json:"short_url"`	}{ShortURL: shortURL_}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("----------------------------------------")
	id:=r.URL.Path[len("/redirect/"):]
	url,err:=getURL(id)
	if err!=nil{
		http.Error(w,"Invalid request",http.StatusNotFound)
	}
	http.Redirect(w,r,url.OriginalURL,http.StatusFound)
}

func getURL(id string) (URL,error){
	url,ok:=urlDB[id];
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url,nil
}
func main(){
	//fmt.Println("Starting URL Shortner")
	http.HandleFunc("/",rootURLHandler)
	http.HandleFunc("/shorten",shortURLHandler)
	http.HandleFunc("/redirect/",redirectURLHandler)
	//generateShortURL("http://example.com")
	fmt.Println("Starting Server on port 8080")
	err:=http.ListenAndServe(":8080",nil)
	if err!=nil {
		fmt.Println("Error on starting server: ",err)
	}
}