package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)


func ConnectDB(ctx context.Context, connStr string) *sql.DB { //?menerima sebua string dan mengembalikan sql.DB
	{
		db, err := sql.Open("postgres", connStr) //?membuka koneksi ke database, mereturn db, err
		if err != nil {
			panic(err)
		}

		err = db.PingContext(ctx)
		if err != nil {
			panic(err)
		}

		return db
	}
}
