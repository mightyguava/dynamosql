-- Partition key required
SELECT * FROM gamescores
-- Partition key required
SELECT * FROM gamescores WHERE Wins = 3
-- $ numbered placeholder is not allowed
SELECT * FROM gamescores WHERE UserId = "101" AND Wins = $1
-- OR may not appear at the top level
SELECT * FROM gamescores WHERE UserId = :UserId OR Wins = 3
-- Partition key may not appear in a nested expression
SELECT * FROM gamescores WHERE UserId = :UserId AND (Wins = 3 OR UserId = "105")
-- Partition key may not appear in a NOT expression
SELECT * FROM gamescores WHERE NOT UserId = :UserId
-- Partition key must be in an equality condition
SELECT * FROM gamescores WHERE UserId > :UserId
-- Partition key must be in an equality condition
SELECT * FROM gamescores WHERE begins_with(UserId, 5)
-- Partition key may not appear twice
SELECT * FROM gamescores WHERE UserId = :UserId AND UserId = :UserId2
-- Sort key may not appear twice
SELECT * FROM gamescores WHERE UserId = :UserId AND GameTitle = "A" AND GameTitle = "B"
-- Sort key only be used with begins_with
SELECT * FROM gamescores WHERE UserId = :UserId AND attribute_exists(GameTitle)
