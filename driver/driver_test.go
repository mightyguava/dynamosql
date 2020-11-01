package driver

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

	"github.com/mightyguava/dynamosql"
	"github.com/mightyguava/dynamosql/testing/fixtures"
)

func TestDriverBind(t *testing.T) {
	sess := fixtures.SetUp(t, fixtures.GameScores)

	driver, err := New(Config{Session: sess}).OpenConnector("")
	require.NoError(t, err)
	db := sql.OpenDB(driver)
	err = db.Ping()
	require.NoError(t, err)

	rows, err := db.Query(`SELECT GameTitle, TopScore FROM gamescores WHERE UserId = "101" AND GameTitle > "Meteor"`)
	require.NoError(t, err)

	var scores []fixtures.GameScore
	for rows.Next() {
		var s fixtures.GameScore
		err = rows.Scan(&s.GameTitle, &s.TopScore)
		require.NoError(t, err)
		scores = append(scores, s)
	}

	rows, err = db.Query(`SELECT GameTitle, TopScore FROM gamescores WHERE UserId = :UserId AND GameTitle > :GameTitle`,
		sql.Named("UserId", "101"),
		sql.Named("GameTitle", "Meteor"))
	require.NoError(t, err)

	for rows.Next() {
		var s fixtures.GameScore
		err = rows.Scan(&s.GameTitle, &s.TopScore)
		require.NoError(t, err)
		scores = append(scores, s)
	}
	repr.Println(scores)
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
					err = rows.Scan(dynamosql.Document(doc))
					require.NoError(t, err, query)
					results = append(results, repr.String(doc))
				}
			} else {
				row := make([]string, len(cols))
				scanRow := make([]interface{}, len(cols))
				for i := range row {
					scanRow[i] = &row[i]
				}
				for rows.Next() {
					err = rows.Scan(scanRow...)
					require.NoError(t, err, query)
					results = append(results, strings.Join(row, ","))
				}
			}
			require.NoError(t, rows.Err())

			result := output{
				Query:   query,
				Results: results,
			}
			g.AssertJson(t, fmt.Sprintf("queries.%02d", i), result)
			i++
		})
	}
}
