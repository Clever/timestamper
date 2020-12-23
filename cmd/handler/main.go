package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"gopkg.in/Clever/kayvee-go.v6/logger"
)

// generate kv config bytes for setting up log routing
//go:generate go-bindata -pkg $GOPACKAGE -o bindata.go kvconfig.yml
//go:generate gofmt -w bindata.go

// Handler encapsulates the external dependencies of the lambda function.
// The example here demonstrates the case where the handler logic involves communicating with S3.
type Handler struct {
	s3API s3iface.S3API
}

// generate a mocks of dependencies for use during testing
//go:generate mockgen -package main -source $PWD/vendor/github.com/aws/aws-sdk-go/service/s3/s3iface/interface.go -destination s3api_mocks_test.go S3API

// Input to the Lambda Function is JSON unmarshal-ed into this struct.
// If you are subscribing to events, then instead of this you should use
// a type from the github.com/aws/aws-lambda-go/events package.
type Input struct {
	Foo string `json:"foo"`
}

// Handle is invoked by the Lambda runtime with the contents of the function input.
func (h Handler) Handle(ctx context.Context, input Input) error {
	// create a request-specific logger, attach it to ctx, and add the Lambda request ID.
	ctx = logger.NewContext(ctx, logger.New(os.Getenv("APP_NAME")))
	if lambdaContext, ok := lambdacontext.FromContext(ctx); ok {
		logger.FromContext(ctx).AddContext("aws-request-id", lambdaContext.AwsRequestID)
	}
	logger.FromContext(ctx).InfoD("received", logger.M{
		"foo": input.Foo,
	})

	if _, err := h.s3API.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Key:    aws.String(input.Foo),
		Bucket: aws.String("bar"),
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := logger.SetGlobalRoutingFromBytes(MustAsset("kvconfig.yml")); err != nil {
		log.Fatalf("Error setting kvconfig: %s", err)
	}

	handler := Handler{
		s3API: s3.New(session.New()),
	}

	if os.Getenv("IS_LOCAL") == "true" {
		// Update input as needed to debug
		input := Input{}
		log.Printf("Running locally with this input: %+v\n", input)
		err := handler.Handle(context.TODO(), input)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		lambda.Start(handler.Handle)
	}
}
