SELECT * FROM gamescores WHERE UserId = :UserId
SELECT * FROM gamescores WHERE UserId = :UserId AND TopScore > :MinTopScore
SELECT * FROM gamescores WHERE UserId = :UserId AND GameTitle BETWEEN :MinGameTitle AND :MaxGameTitle AND TopScore > :MinTopScore
SELECT * FROM gamescores WHERE UserId = "103" AND GameTitle BETWEEN "Galaxy" AND "Meteor" AND TopScore > 1000
SELECT * FROM gamescores WHERE UserId = "103" AND begins_with(GameTitle, "Galaxy")
-- Column projections
SELECT UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3] FROM gamescores WHERE UserId = :UserId
-- Document projection
SELECT document(UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3]) FROM gamescores WHERE UserId = :UserId
-- Reserved word substitution
SELECT title, year FROM movies WHERE title = :title AND year > 2009 AND escaped = TRUE
-- Fields that have dots in them. {"foo.bar": "a"} and {"foo": {"bar": "b"}} are different.
SELECT `foo.bar`, `foo`.`bar` FROM movies WHERE title = :title
-- Global Secondary Index with different hash key
SELECT * FROM gamescores USE INDEX (GameTitleIndex) WHERE GameTitle = :title AND UserId > "45"
-- LIMIT
SELECT * FROM gamescores WHERE UserId = "103" LIMIT 1
