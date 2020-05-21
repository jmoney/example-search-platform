package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jmoney8080/example-search-platform/indexer/internal/handle"
)

func main() {
	lambda.Start(handle.Request)
}
