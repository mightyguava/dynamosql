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
