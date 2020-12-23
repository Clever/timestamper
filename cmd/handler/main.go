package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"gopkg.in/Clever/kayvee-go.v6/logger"
)

// generate kv config bytes for setting up log routing
//go:generate go-bindata -pkg $GOPACKAGE -o bindata.go kvconfig.yml
//go:generate gofmt -w bindata.go

// Handler encapsulates the external dependencies of the lambda function.
type Handler struct {
}

// Input to the Lambda Function is JSON unmarshal-ed into this struct.
// If you are subscribing to events, then instead of this you should use
// a type from the github.com/aws/aws-lambda-go/events package.
type Input struct {
}

// Output ...
type Output struct {
	Timestamp time.Time `json:"timestamp"`
}

// Handle is invoked by the Lambda runtime with the contents of the function input.
func (h Handler) Handle(ctx context.Context, input Input) (Output, error) {
	// create a request-specific logger, attach it to ctx, and add the Lambda request ID.
	ctx = logger.NewContext(ctx, logger.New(os.Getenv("APP_NAME")))
	if lambdaContext, ok := lambdacontext.FromContext(ctx); ok {
		logger.FromContext(ctx).AddContext("aws-request-id", lambdaContext.AwsRequestID)
	}
	logger.FromContext(ctx).InfoD("received", logger.M{})

	now := time.Now().UTC()

	return Output{Timestamp: now}, nil
}

func main() {
	if err := logger.SetGlobalRoutingFromBytes(MustAsset("kvconfig.yml")); err != nil {
		log.Fatalf("Error setting kvconfig: %s", err)
	}

	handler := Handler{}

	if os.Getenv("IS_LOCAL") == "true" {
		// Update input as needed to debug
		input := Input{}
		log.Printf("Running locally with this input: %+v\n", input)
		output, err := handler.Handle(context.TODO(), input)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	} else {
		lambda.Start(handler.Handle)
	}
}
