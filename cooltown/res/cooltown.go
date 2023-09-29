package cooltown

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Song struct {
	Audio string
}

func findTrack(w http.ResponseWriter, r *http.Request) {
	const search = "http://localhost:3001/search"
	songInfo, err := http.Post(search, "text/plain", r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if songInfo.StatusCode == 400 {
		w.WriteHeader(400)
		return
	} else if songInfo.StatusCode == 404 {
		w.WriteHeader(404)
		return
	} else if songInfo.StatusCode == 500 {
		w.WriteHeader(500)
		return
	}

	title, err := io.ReadAll(songInfo.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var m = make(map[string]string)
	err = json.Unmarshal(title, &m)
	if err != nil {
		w.WriteHeader(500) // unmarshal failed
		return
	}

	extractedTitle := m["Id"]
	extractedTitle = strings.Replace(extractedTitle, " ", "+", -1)
	// now that we have the title must return the full song
	getAudio := "http://localhost:3000/tracks/" + string(extractedTitle)
	fullAudioResponse, err := http.Get(getAudio)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if fullAudioResponse.StatusCode != 200 {
		if fullAudioResponse.StatusCode == 400 {
			w.WriteHeader(400)
			return
		} else if fullAudioResponse.StatusCode == 404 {
			w.WriteHeader(404)
			return
		} else {
			w.WriteHeader(500)
			return
		}
	}
	audioBytes, err := io.ReadAll(fullAudioResponse.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	m = make(map[string]string)
	err = json.Unmarshal(audioBytes, &m)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	songAudio := m["Audio"]

	if songAudio == "" {
		w.WriteHeader(400) // empty song
		return
	}
	songOut := Song{songAudio}

	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(songOut)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Store */
	r.HandleFunc("/cooltown", findTrack).Methods("POST")
	return r
}
