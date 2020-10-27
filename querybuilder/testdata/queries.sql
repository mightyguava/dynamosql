SELECT * FROM gamescores WHERE UserId = :UserId
SELECT * FROM gamescores WHERE UserId = :UserId AND TopScore > :MinTopScore
SELECT * FROM gamescores WHERE UserId = :UserId AND GameTitle BETWEEN "Galaxy" AND "Meteor" AND TopScore > :MinTopScore