SELECT * FROM gamescores WHERE UserId = :UserId
SELECT * FROM gamescores WHERE UserId = :UserId AND TopScore > :MinTopScore
SELECT * FROM gamescores WHERE UserId = :UserId AND GameTitle BETWEEN :MinGameTitle AND :MaxGameTitle AND TopScore > :MinTopScore
SELECT * FROM gamescores WHERE UserId = "103" AND GameTitle BETWEEN "Galaxy" AND "Meteor" AND TopScore > 1000
