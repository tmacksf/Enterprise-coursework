package tracks

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"tracks/repository"
)

func updateTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var t repository.Track

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if id != t.Id {
		w.WriteHeader(400)
		return
	}

	if t.Id == "" || t.Audio == "" {
		w.WriteHeader(400)
		return
	}

	_, doesTrackExist := repository.ReadTrack(id)
	if doesTrackExist < 0 {
		w.WriteHeader(500)
		return
	}

	if doesTrackExist == 0 {
		// track does not exist
		created := repository.CreateTrack(t)
		if created > 0 {
			w.WriteHeader(201) // track was created
			return
		} else {
			w.WriteHeader(500)
			return
		}
	}

	if doesTrackExist > 0 {
		// track does exist
		updated := repository.UpdateTrack(t)
		if updated == 1 {
			w.WriteHeader(204) // track was updated
			return
		} else {
			w.WriteHeader(500)
			return
		}
	}
}

func readTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if t, n := repository.ReadTrack(id); n > 0 {
		d := repository.Track{Id: t.Id, Audio: t.Audio}
		w.WriteHeader(200) /* OK */
		json.NewEncoder(w).Encode(d)
	} else if n == 0 {
		w.WriteHeader(404) /* Not Found */
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func listTracks(w http.ResponseWriter, r *http.Request) {
	// TODO figure out of n has to be greater than 0
	if t, n := repository.ListTracks(); n > 0 {
		w.WriteHeader(200) /* OK */
		json.NewEncoder(w).Encode(t)
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if num := repository.DeleteTrack(id); num > 0 {
		w.WriteHeader(204)
	} else if num == 0 {
		w.WriteHeader(400) // No track with this ID
	} else {
		w.WriteHeader(500) // Error
	}
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Store */
	r.HandleFunc("/tracks/{id}", updateTrack).Methods("PUT")
	r.HandleFunc("/tracks/{id}", deleteTrack).Methods("DELETE")
	/* Document */
	r.HandleFunc("/tracks/{id}", readTrack).Methods("GET")
	r.HandleFunc("/tracks", listTracks).Methods("GET")
	return r
}
