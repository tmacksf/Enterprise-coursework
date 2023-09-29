package search

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
)

type output struct {
	Id string
}

type apiResponse struct {
	Status    string          `json:"status"`
	SongInfo  json.RawMessage `json:"result"`
	ErrorInfo json.RawMessage `json:"error"`
}

type errorInfo struct {
	ErrorCode int `json:"error_code"`
}

func Search(w http.ResponseWriter, r *http.Request) {
	audioBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	var m = make(map[string]string)
	err = json.Unmarshal(audioBytes, &m)
	if err != nil {
		w.WriteHeader(400) // unmarshal failed
		return
	}

	audio := m["Audio"]
	if audio == "" {
		w.WriteHeader(400) // bad request (no audio)
		return
	}
	if id, n := searchAPI(audio); n == 1 {
		w.WriteHeader(200) // OK
		song := output{id}
		err = json.NewEncoder(w).Encode(song)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	} else if n == 400 {
		w.WriteHeader(400) // Not Found
	} else if n == 404 {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(500) // todo figure it out
	}
}

func searchAPI(Audio string) (string, int64) {
	const apikey = "7d4b2a59f3f3fa092762a11fb25f7d2b"
	const link = "https://api.audd.io/"

	if apikey == "" {
		err := fmt.Errorf("no api key")
		fmt.Println(err)
		return "", 500
	}

	form := url.Values{}
	form.Add("api_token", apikey)
	form.Add("audio", Audio)

	response, err := http.PostForm(link, form)
	if err != nil {
		return "", 500
	}

	if response.StatusCode != http.StatusOK {
		return "", 400 // todo figure out if this should be a 500 or a 400
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", 500 // cannot decode the response body = 500?
	}

	var responseInfo apiResponse
	err = json.Unmarshal(body, &responseInfo)
	if err != nil {
		return "", 500
	}

	if responseInfo.Status != "success" {
		// error processing (if an api token is wrong then 400. If the song is not right then 404, else 500)
		//errorInfo := make(map[string]string)
		var errInfo errorInfo
		if err = json.Unmarshal(responseInfo.ErrorInfo, &errInfo); err != nil {
			return "", 500
		}

		if errInfo.ErrorCode == 900 {
			return "", 400 // bad request because of bad api key
		} else if errInfo.ErrorCode == 300 {
			return "", 404 // song not found
		} else if errInfo.ErrorCode == 700 {
			return "", 404 // song not found
		} else {
			return "", 500 // not sure if this should be 400 or 500
		}
	}

	songInfo := make(map[string]string)
	err = json.Unmarshal(responseInfo.SongInfo, &songInfo)
	if err != nil {
		return "", 500
	}
	if songInfo == nil {
		return "", 404
	}
	songTitle := songInfo["title"]
	//songTitle = strings.Replace(songTitle, " ", "+", -1)
	return songTitle, 1
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Store */
	r.HandleFunc("/search", Search).Methods("POST")
	return r
}
