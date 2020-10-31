SELECT * FROM gamescores WHERE UserId = :UserId
SELECT * FROM gamescores WHERE UserId = :UserId AND TopScore > :MinTopScore
SELECT * FROM gamescores WHERE UserId = :UserId AND GameTitle BETWEEN :MinGameTitle AND :MaxGameTitle AND TopScore > :MinTopScore
SELECT * FROM gamescores WHERE UserId = "103" AND GameTitle BETWEEN "Galaxy" AND "Meteor" AND TopScore > 1000
SELECT * FROM gamescores WHERE UserId = "103" AND begins_with(GameTitle, "Galaxy")
SELECT UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3] FROM gamescores WHERE UserId = :UserId
