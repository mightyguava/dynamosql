# dynamosql - DynamoDB SQL Driver

![Test Status](https://github.com/mightyguava/dynamosql/workflows/Test/badge.svg)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/mightyguava/dynamosql)

`dynamosql` is a complete Go SQL driver for DynamoDB. It transforms SQL into DynamoDB requests like `Query`, `Scan`, `PutItem`, `UpdateItem`, `DeleteItem`, and maps the results back to SQL. It

* Makes working with DynamoDB much more pleasant.
* Transparently generate KeyConditionExpression, ConditionExpression, etc
* Supports advanced SQL driver features like named parameters, and slice/map parameters.
* Supports marshaling and unmarshaling using [`dynamodbattribute`](https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/).
* Works with any library that supports `database/sql`.
* Does not magically add any features like JOIN or cross partition queries.

There is also a CLI in `cmd/dynamosql` that can be used as a commandline interface to DynamoDB. For authentication, it accepts the same environment variables, config, and flags as the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

## Inspiration

SQL is a great query language, but traditional SQL databases tend not to scale so well. DynamoDB is a great scalable database, but it's query language is lacking. So why not combine them? Hence the inspiration for `dynamosql`.

## Usage

### CLI

```bash
go get github.com/mightyguava/dynamosql/cmd/dynamosql
aws configure # if you need to
dynamosql
> SELECT * FROM my_table WHERE hash_key = 123
```

### Driver

There are 3 ways to open a connection.

via `database/sql`:

```go
db, err := sql.Open("dynamodb://?region=us-west-2")
```
passing a `Session` into the driver

```go
sess, err := session.New(&aws.Config{
	Region: aws.String("us-west-2"),
})
db := dynamosql.NewWithSession(sess)
```

passing a DynamoDB client into the driver

```
ddb := dynamodb.New(sess)
db := dynamosql.NewWithClient(ddb)
```


### Permissions

`dynamosql` requires the following permissions to be granted.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "DynamoDBTableAccess",
            "Effect": "Allow",
            "Action": [
                "dynamodb:BatchGetItem",
                "dynamodb:BatchWriteItem",
                "dynamodb:ConditionCheckItem",
                "dynamodb:PutItem",
                "dynamodb:DescribeTable",
                "dynamodb:DeleteItem",
                "dynamodb:GetItem",
                "dynamodb:Scan",
                "dynamodb:Query",
                "dynamodb:UpdateItem"
            ],
            "Resource": "arn:aws:dynamodb:us-west-2:123456789012:table/TableName"
        }
    ]
}
```

## SQL Mappings

| SQL | DynamoDB | Notes |
| --- | --- | --- |
| SELECT | Query |
| (TODO) SCAN | Scan | Optionally turn SELECT into SCAN when no keys are present in WHERE
| INSERT | PutItem/TransactWriteItem | Errors if key exists. Uses TransactWriteItem to insert up to 25 items |
| REPLACE ... RETURNING | PutItem/TransactWriteItem | Overwrites existing document.  Uses TransactWriteItem to insert up to 25 items |
| (TODO) UPDATE | UpdateItem | |
| CREATE TABLE | CreateTable | supports global and local secondary indexes |
| (TODO) ALTER TABLE | | |

## Example

A fairly complete example of driver usage. Error checking omitted for brevity.

```go
sess := session.Must(session.New(aws.Config{Region: aws.String()"us-west-2")}))
db := dynamosql.NewWithSession(sess)
_, err := db.Exec(`
CREATE TABLE movies (
  title string HASH KEY,
  year number RANGE KEY,
  director string,
  actors []string,

  LOCAL INDEX title_director RANGE(director) PROJECTION ALL,
  GLOBAL INDEX director_year HASH(director) RANGE(year)
      PROJECTION INCLUDE (title, actors)
      PROVISIONED THROUGHPUT READ 1 WRITE 1
)
PROVISIONED THROUGHPUT READ 1 WRITE 1;
`)
type Movie struct {
    title string
    year int
}
movies := []movie{
    {"Rush Hour", 1994},
    {"Die Hard", 1988},
}
_, err := db.Exec(`INSERT INTO movies VALUES (?)`, movies)
var rushHour Movie
row := db.QueryRow(`SELECT * FROM movies WHERE title = :name`, sql.Named("name", "Rush Hour"))
err := row.Scan(dynamosql.Document(&rushHour))
fmt.Println(rushHour)
```

## Grammar

```
AST = (Select | InsertOrReplace | CreateTable | DropTable) ";"? .

Select = "SELECT" ProjectionExpression "FROM" <field> ("USE" "INDEX" "(" <ident> ")")? ("WHERE" AndExpression)? ("ASC" | "DESC")? ("LIMIT" <number>)? .
ProjectionExpression = ("*" | ("document" "(" "*" ")")) | (ProjectionColumn ("," ProjectionColumn)*) .
ProjectionColumn = FunctionExpression | DocumentPath .
FunctionExpression = <ident> "(" FunctionArgument ("," FunctionArgument)* ")" .
FunctionArgument = DocumentPath | Value .
DocumentPath = PathFragment ("." PathFragment)* .
PathFragment = <field> ("[" <number> "]")* .
Value = <number> | <string> | <bool> | <null> | (":" <ident>) | "?" .
AndExpression = Condition ("AND" Condition)* .
Condition = ("(" ConditionExpression ")") | ("NOT" NotCondition) | ConditionOperand | FunctionExpression .
ConditionExpression = AndExpression ("OR" AndExpression)* .
NotCondition = Condition .
ConditionOperand = DocumentPath ConditionRHS .
ConditionRHS = Compare | ("BETWEEN" Between) | ("IN" "(" In ")") .
Compare = ("<>" | "<=" | ">=" | "=" | "<" | ">" | "!=") Operand .
Operand = Value | DocumentPath .
Between = Operand "AND" Operand .
In = Value ("," Value)* .

InsertOrReplace = ("INSERT" | "REPLACE") "INTO" <field> "VALUES" "(" InsertTerminal ")" ("," "(" InsertTerminal ")")* ("RETURNING" ("NONE" | "ALL_OLD"))? .
InsertTerminal = <number> | <string> | <bool> | <null> | (":" <ident>) | "?" | JSONObject .
JSONObject = "{" (JSONObjectEntry ("," JSONObjectEntry)* ","?)? "}" .
JSONObjectEntry = (<ident> | <string>) ":" JSONValue .
JSONValue = <number> | <string> | <bool> | <null> | JSONObject | JSONArray .
JSONArray = "[" (JSONValue ("," JSONValue)* ","?)? "]" .

CreateTable = "CREATE" "TABLE" <field> "(" CreateTableEntry ("," CreateTableEntry)* ")" ProvisionedThroughput .
CreateTableEntry = GlobalSecondaryIndex | LocalSecondaryIndex | TableAttr .
GlobalSecondaryIndex = "GLOBAL" "SECONDARY"? "INDEX" <field> "HASH" "(" <field> ")" ("RANGE" "(" <field> ")")? "PROJECTION" Projection ProvisionedThroughput .
Projection = "KEYS_ONLY" | "ALL" | ("INCLUDE" (<field> ("," <field>)*)) .
ProvisionedThroughput = "PROVISIONED" "THROUGHPUT" "READ" <number> "WRITE" <number> .
LocalSecondaryIndex = "LOCAL" "SECONDARY"? "INDEX" <field> "RANGE" "(" <field> ")" "PROJECTION" Projection .
TableAttr = <field> <type> (("HASH" | "RANGE") "KEY")? .

DropTable = "DROP" "TABLE" <field> .
```
