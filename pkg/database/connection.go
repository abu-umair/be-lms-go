package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context, connStr string) *sqlx.DB { //?menerima sebua string dan mengembalikan sql.DB
	{
		db, err := sqlx.Open("postgres", connStr) //?membuka koneksi ke database, mereturn db, err
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
