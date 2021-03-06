SELECT * FROM movies
SELECT title, year FROM movies
SELECT title, year FROM movies WHERE title = "The Dark Knight"
SELECT title, year FROM movies WHERE title = "The Dark Knight" AND year >= 2009
SELECT title, year FROM movies WHERE title = "The Dark Knight" AND year BETWEEN 2009 AND 2015
SELECT title, year FROM movies WHERE title = "The Dark Knight" AND year BETWEEN 2009 AND 2015 AND actor = "Will Smith"
SELECT title, year FROM movies WHERE title = "The Dark Knight" AND (year BETWEEN 2009 AND 2015 OR actor = "Will Smith")
SELECT * FROM movies WHERE title = :title
SELECT * FROM movies WHERE title = :title AND attribute_exists(year)
SELECT * FROM movies WHERE title = :title AND begins_with(actor, "Will")
-- Projection tests https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.Attributes.html
-- Standard projection
SELECT UserId, TopScore FROM gamescores WHERE UserId = :UserId
-- List projections
SELECT Scores[3], Scores[3][2] FROM gamescores WHERE UserId = :UserId
-- Map element
SELECT Studio.Name, Studio.Name.FirstName, Studio.Employees[3] FROM gamescores WHERE UserId = :UserId
-- All the projections
SELECT UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3] FROM gamescores WHERE UserId = :UserId
-- Project fields as document
SELECT document(UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3]) FROM gamescores WHERE UserId = :UserId
-- Project fields and document
SELECT UserId, document(TopScore) FROM gamescores WHERE UserId = :UserId
-- Quoted keywords
SELECT `SELECT`.`foo` FROM gamescores WHERE UserId = :UserId
-- document(*) is treated the same as *
SELECT document(*) FROM gamescores WHERE UserId = :UserId
-- Use an index
SELECT * FROM movies USE INDEX (some_index) WHERE UserId = :UserId
-- Boolean parsing
SELECT * FROM movies WHERE UserId = True
-- DESC parsing
SELECT * FROM movies WHERE UserId = :UserId DESC
-- Dot in table name
SELECT * FROM namespaced.movies WHERE UserId = :UserId
-- positional placeholders (?)
SELECT * FROM gamescores WHERE UserId = ? AND TopScore > ?
-- create table
CREATE TABLE movies (title STRING, year NUMBER) PROVISIONED THROUGHPUT READ 1 WRITE 1;
CREATE TABLE movies (title STRING HASH KEY, year NUMBER RANGE KEY) PROVISIONED THROUGHPUT READ 1 WRITE 1;
CREATE TABLE movies (title STRING, year NUMBER, GLOBAL SECONDARY INDEX year_title HASH(year) RANGE(title) PROJECTION ALL PROVISIONED THROUGHPUT READ 1 WRITE 1) PROVISIONED THROUGHPUT READ 1 WRITE 1;
CREATE TABLE movies (title STRING, year NUMBER, LOCAL SECONDARY INDEX year_index RANGE(year) PROJECTION ALL) PROVISIONED THROUGHPUT READ 1 WRITE 1;
-- drop table
DROP TABLE movies;