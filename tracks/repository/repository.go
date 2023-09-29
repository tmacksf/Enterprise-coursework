package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Repository struct {
	DB *sql.DB
}

var repo Repository

/* Functions that deal with creation of database (done in the labs) */

func Init() {
	if db, err := sql.Open("sqlite3", "./tmp/test.db"); err == nil {
		repo = Repository{DB: db}
	} else {
		log.Fatal("Database initialisation")
	}
}

func Create() int {
	const sqlStmt = "CREATE TABLE IF NOT EXISTS Tracks" +
		"(Id TEXT PRIMARY KEY, Audio TEXT)"
	if _, err := repo.DB.Exec(sqlStmt); err == nil {
		return 0
	} else {
		return -1
	}
}

func Clear() int {
	const sqlStmt = "DELETE FROM Tracks"
	if _, err := repo.DB.Exec(sqlStmt); err == nil {
		return 0
	} else {
		return -1
	}
}

/***Functions that deal with the coursework***/

func ReadTrack(Id string) (Track, int) {
	sqlStmt := "SELECT * FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sqlStmt); err == nil {
		defer stmt.Close()
		var t Track
		row := stmt.QueryRow(Id)
		if err := row.Scan(&t.Id, &t.Audio); err == nil {
			return t, 1
		} else {
			return t, 0
		}
	} else {
		return Track{}, -1
	}
}

func UpdateTrack(t Track) int64 {
	const sqlStmt = " UPDATE Tracks SET Audio = ? WHERE id = ?"
	if stmt, err := repo.DB.Prepare(sqlStmt); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Audio, t.Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

func CreateTrack(t Track) int64 {
	sqlStmt := "INSERT INTO Tracks(Id, Audio) " + "VALUES (?, ?)"
	if stmt, err := repo.DB.Prepare(sqlStmt); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Id, t.Audio); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

func ListTracks() ([]string, int64) {
	sqlStmt := "SELECT Id FROM Tracks"
	var tracks []string
	if rows, err := repo.DB.Query(sqlStmt); err == nil {
		defer rows.Close()
		for rows.Next() {
			var Id string
			if err := rows.Scan(&Id); err == nil {
				tracks = append(tracks, Id)
			}
		}
		return tracks, 1
	}
	return tracks, -1
}

func DeleteTrack(Id string) int64 {
	sqlStmt := "DELETE FROM Tracks WHERE id = ?"
	if stmt, err := repo.DB.Prepare(sqlStmt); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}
