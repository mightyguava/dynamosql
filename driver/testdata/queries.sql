SELECT * FROM gamescores WHERE UserId = "101"
SELECT * FROM gamescores WHERE UserId = "101" AND TopScore > 30
SELECT * FROM gamescores WHERE UserId = "103" AND GameTitle BETWEEN "Galaxy" AND "Meteor" AND TopScore > 1000
SELECT GameTitle, TopScore FROM gamescores WHERE UserId = "101"
SELECT title, year FROM movies WHERE title = "World War Z"
SELECT title, info.rating FROM movies WHERE title = "World War Z"
