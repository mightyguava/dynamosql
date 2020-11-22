package dynamosql

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/testing/fixtures"
	"github.com/mightyguava/dynamosql/testing/testutil"
)

func TestDriverBind(t *testing.T) {
	sess := fixtures.SetUp(t, fixtures.GameScores)

	driver, err := New(Config{Session: sess}).OpenConnector("")
	require.NoError(t, err)
	db := sql.OpenDB(driver)
	err = db.Ping()
	require.NoError(t, err)

	readRows := func(rows *sql.Rows) []fixtures.GameScore {
		var scores []fixtures.GameScore
		for rows.Next() {
			var s fixtures.GameScore
			err = rows.Scan(&s.GameTitle, &s.TopScore)
			require.NoError(t, err)
			scores = append(scores, s)
		}
		require.NoError(t, rows.Err())
		return scores
	}
	expected := []fixtures.GameScore{
		{GameTitle: "Meteor Blasters", TopScore: 1000},
		{GameTitle: "Starship X", TopScore: 24}}

	t.Run("fixed params", func(t *testing.T) {
		rows, err := db.Query(`SELECT GameTitle, TopScore FROM gamescores WHERE UserId = "101" AND GameTitle > "Meteor"`)
		require.NoError(t, err)
		require.Equal(t, expected, readRows(rows))
	})

	t.Run("named params", func(t *testing.T) {
		rows, err := db.Query(`SELECT GameTitle, TopScore FROM gamescores WHERE UserId = :UserId AND GameTitle > :GameTitle`,
			sql.Named("UserId", "101"),
			sql.Named("GameTitle", "Meteor"))
		require.NoError(t, err)
		require.Equal(t, expected, readRows(rows))
	})

	t.Run("positional params", func(t *testing.T) {
		rows, err := db.Query(`SELECT GameTitle, TopScore FROM gamescores WHERE UserId = ? AND GameTitle > ?`, "101", "Meteor")
		require.NoError(t, err)
		require.Equal(t, expected, readRows(rows))
	})
}

func TestDriverGolden(t *testing.T) {
	sess := fixtures.SetUp(t, fixtures.GameScores, fixtures.Movies)

	driver, err := New(Config{Session: sess}).OpenConnector("")
	require.NoError(t, err)
	db := sql.OpenDB(driver)
	err = db.Ping()
	require.NoError(t, err)

	type output struct {
		Query   string
		Results []interface{}
	}

	queries, err := os.Open("testdata/queries.sql")
	require.NoError(t, err)
	defer queries.Close()
	scanner := bufio.NewScanner(queries)
	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithNameSuffix(".golden.json"))
	i := 0
	for scanner.Scan() {
		query := scanner.Text()
		if strings.HasPrefix(query, "--") {
			// skip comments
			continue
		}
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rows, err := db.Query(query)
			require.NoError(t, err, query)

			var results []interface{}
			cols, err := rows.Columns()
			require.NoError(t, err, query)
			results = append(results, strings.Join(cols, ","))

			if cols[0] == "document" {
				var doc interface{}
				if strings.Contains(query, "FROM movies") {
					doc = &fixtures.Movie{}
				} else if strings.Contains(query, "FROM gamescores") {
					doc = &fixtures.GameScore{}
				} else {
					panic("unexpected code path")
				}
				for rows.Next() {
					err = rows.Scan(Document(doc))
					require.NoError(t, err, query)
					results = append(results, repr.String(doc))
				}
			} else {
				row := make([]sql.NullString, len(cols))
				scanRow := make([]interface{}, len(cols))
				for i := range row {
					scanRow[i] = &row[i]
				}
				for rows.Next() {
					err = rows.Scan(scanRow...)
					require.NoError(t, err, query)
					results = append(results, strings.Join(stringSlice(row), ","))
				}
			}
			require.NoError(t, rows.Err())

			result := output{
				Query:   query,
				Results: results,
			}
			g.Assert(t, fmt.Sprintf("queries.%02d", i), []byte(testutil.MarshalJSON(result)))
			i++
		})
	}
}

func stringSlice(ns []sql.NullString) []string {
	ss := make([]string, len(ns))
	for i := range ns {
		if !ns[i].Valid {
			ss[i] = "NULL"
		} else {
			ss[i] = ns[i].String
		}
	}
	return ss
}

func TestInsertAndQuery(t *testing.T) {
	fix := fixtures.Movies
	fix.Data = nil
	sess := fixtures.SetUp(t, fix)

	driver, err := New(Config{Session: sess}).OpenConnector("")
	require.NoError(t, err)
	db := sql.OpenDB(driver)
	err = db.Ping()
	require.NoError(t, err)

	prisoners := fixtures.Movie{"Prisoners", 2013, fixtures.MovieInfo{}}
	rushHour := fixtures.Movie{"Rush Hour", 1998, fixtures.MovieInfo{}}
	forrestGump := fixtures.Movie{"Forrest Gump", 1994, fixtures.MovieInfo{}}
	inception := fixtures.Movie{"Inception", 2010, fixtures.MovieInfo{}}
	dieHard := fixtures.Movie{"Die Hard", 1988, fixtures.MovieInfo{}}

	v, err := db.Exec("INSERT INTO movies VALUES (?)", []fixtures.Movie{
		prisoners, rushHour,
	})
	require.NoError(t, err)
	rows, err := v.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(2), rows)

	v, err = db.Exec("INSERT INTO movies VALUES (?)", forrestGump)
	require.NoError(t, err)
	rows, err = v.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)

	v, err = db.Exec(`
INSERT INTO movies VALUES 
('{"title":"Inception", "year": 2010}')
`)
	require.NoError(t, err)
	v, err = db.Exec(`
INSERT INTO movies VALUES
({title:"Die Hard", year:1988})
`)
	require.NoError(t, err)
	rows, err = v.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)

	for _, m := range []fixtures.Movie{prisoners, rushHour, forrestGump, inception, dieHard} {
		row := db.QueryRow(`SELECT * FROM movies WHERE title = ?`, m.Title)
		var movie fixtures.Movie
		require.NoError(t, row.Scan(Document(&movie)))
		require.Equal(t, m, movie)
	}
}
