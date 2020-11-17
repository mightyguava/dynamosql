module github.com/mightyguava/dynamosql/cmd/dynamosql

go 1.15

require (
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/aws/aws-sdk-go v1.28.7
	github.com/mightyguava/dynamosql v0.0.0-20201106003410-0b8f4855fb23
	github.com/xo/usql v0.7.8
)

replace github.com/mightyguava/dynamosql => ../../
