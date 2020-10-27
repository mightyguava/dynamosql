package driver

import (
	"database/sql"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/testing/fixtures"
)

func TestDriver(t *testing.T) {
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
