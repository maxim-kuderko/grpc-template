package primary

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"github.com/spf13/viper"
	"runtime"
)

type Db struct {
	client *sql.DB
}

func NewDb(v *viper.Viper) Repo {
	return &Db{
		client: createSqlConnection(v.GetString(`PRIMARY_MYSQL_DSN`)),
	}
}

func createSqlConnection(dsn string) *sql.DB {
	db, err := sql.Open(`mysql`, dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(runtime.GOMAXPROCS(0) * 2)
	db.SetMaxIdleConns(runtime.GOMAXPROCS(0))
	return db
}

const GET_QUERY = "select value from data where `key` = ?"

func (d *Db) Get(r requests.Get) (responses.Get, error) {
	output := responses.Get{}
	err := d.client.QueryRow(GET_QUERY, r.Key).Scan(&output.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			output.StatusCode = 404
			return output, nil
		}
	}
	return output, err
}
