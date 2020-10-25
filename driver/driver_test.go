package driver

import (
	"database/sql"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/testing/fixtures"
)

func TestDriver(t *testing.T) {
	sess := fixtures.SetUp(t, fixtures.Movies)

	driver, err := New(Config{Session: sess}).OpenConnector("")
	require.NoError(t, err)
	db := sql.OpenDB(driver)
	err = db.Ping()
	require.NoError(t, err)

	rows, err := db.Query(`SELECT * FROM movies WHERE title = "The Dark Knight"`)
	require.NoError(t, err)

	var movies []fixtures.Movie
	for rows.Next() {
		var m fixtures.Movie
		err = rows.Scan(&m.Title, &m.Year)
		require.NoError(t, err)
		movies = append(movies, m)
	}
	repr.Println(movies)
}
