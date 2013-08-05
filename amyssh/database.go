package amyssh

import (
	"fmt"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var _ = fmt.Printf // deleteme

type Connection struct {
	db *sql.DB
}

func dbStr(cfg DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
}

func NewCon(cfg Config) *Connection {
	con := Connection{}
	var err error
	con.db, err = sql.Open("mysql", dbStr(cfg.Database))

	if err != nil {
		panic(err)
	}
	return &con
}

func (con *Connection) FetchKeys(hostTags []string, userTags []string) (keys []string) {
	// hostTags := []string{"all", "s-rt", "staging"}

	typeParam := "host"

	// TODO: find better way to use prepared statement escaping
	params := append([]interface{}{typeParam})
	tagPlaceholders := make([]string, len(hostTags))
	for i, tag := range hostTags {
		tagPlaceholders[i] = "?"
		params = append(params, tag)
	}

	query := fmt.Sprintf("SELECT DISTINCT `key` FROM ssh_keys k "+
		"JOIN tags t ON t.type=? AND t.label IN (%s) AND k.idssh_keys = t.idssh_keys",
		strings.Join(tagPlaceholders, ","))

	//TODO investigate append on pre 'made' array of some arbitrary size
	row, err := con.db.Query(query, params...)
	if err != nil {
		panic(err)
	}

	var result []string
	var key string
	for row.Next() {
		row.Scan(&key)
		result = append(result, key)
	}
	// Show()
	return result
}

func Show(v interface{}) {
	fmt.Printf("%+v\n", v)
}
