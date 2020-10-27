module github.com/mightyguava/dynamosql

go 1.15

require (
	github.com/alecthomas/kong v0.2.11
	github.com/alecthomas/participle v0.6.0
	github.com/alecthomas/repr v0.0.0-20201006074542-804e374aceb1
	github.com/aws/aws-sdk-go v1.35.14
	github.com/sebdah/goldie/v2 v2.5.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
)

replace github.com/sebdah/goldie/v2 => github.com/mightyguava/goldie/v2 v2.5.2-0.20201027035444-df66e94788ee
