package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

type database struct {
	db *sql.DB
}

type ShortcutReference struct {
	ShortcutID int
	ProfileIDs []int
}

func (db database) InitDB() error {
	var err error
	db.db, err = sql.Open("sqlite3", "./ppl.db")
	if err != nil {
		l.Logger.Error("Error opening database", "Error", err)
		return err
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS shortcuts (
        shortcutid INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT
    );
    CREATE TABLE IF NOT EXISTS profile_versions (
        versionid INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        major INTEGER NOT NULL,
        minor INTEGER NOT NULL,
        revision INTEGER NOT NULL,
        build INTEGER NOT NULL
    );
    CREATE TABLE IF NOT EXISTS profiles (
        profileid INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        profilename TEXT NOT NULL,
        versionid INTEGER,
        FOREIGN KEY (versionid) REFERENCES profile_versions(versionid)
    );
    CREATE TABLE IF NOT EXISTS shortcut_profiles (
        shortcutid INTEGER,
        profileid INTEGER,
        FOREIGN KEY (shortcutid) REFERENCES shortcuts(shortcutid),
        FOREIGN KEY (profileid) REFERENCES profiles(profileid),
        PRIMARY KEY (shortcutid, profileid)
    );
    `
	_, err = db.db.Exec(sqlStmt)
	if err != nil {
		l.Logger.Error("Error creating tables", "Error", err)
		return err
	}

	return nil
}

func (db database) CloseDB() {
	db.db.Close()
}
