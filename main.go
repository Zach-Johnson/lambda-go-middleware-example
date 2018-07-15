package main

import (
	"context"
	"encoding/json"
	"errors"

	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

// Event is a generic Event that the Lambda is invoked with.
type Event struct {
	GenericField string
}

func main() {
	lambda.Start(
		StayToasty( // Ping event handler - returns if a ping is received.
			ParseEvent( // Parse the standard event - return on error.
				Process, // Process the standard event.
			),
		),
	)
}

// HandlerFunc is a generic JSON Lambda handler used to chain middleware.
type HandlerFunc func(context.Context, json.RawMessage) (interface{}, error)

// MiddlewareFunc is a generic middleware example that takes in a HandlerFunc
// and calls the next middleware in the chain.
func MiddlewareFunc(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, data json.RawMessage) (interface{}, error) {
		return next(ctx, data)
	})
}

// StayToasty is a middleware that checks for a ping and returns if it was a ping
// otherwise it will call the next middleware in the chain.
// Can be used to keep a Lambda function warm.
func StayToasty(next HandlerFunc) HandlerFunc {
	type ping struct {
		Ping string `json:"ping"`
	}

	return HandlerFunc(func(ctx context.Context, data json.RawMessage) (interface{}, error) {
		var p ping
		// If unmarshal into the ping struct is successful and there was a value in ping, return out.
		if err := json.Unmarshal(data, &p); err == nil && p.Ping != "" {
			log.Println("ping")
			return "pong", nil
		}
		// Otherwise it's a regular request, call the next middleware.
		return next(ctx, data)
	})
}

// ParseEvent is a middleware that unmarshals an Event before passing control.
func ParseEvent(h EventHandler) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, data json.RawMessage) (interface{}, error) {
		var event Event

		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}

		if event.GenericField == "" {
			return nil, errors.New("GenericField must be populated")
		}

		return h(&ctx, &event)
	})
}

// EventHandler is the function signature to process the Event.
type EventHandler func(*context.Context, *Event) (interface{}, error)

// Process satisfies EventHandler and processes the Event.
func Process(ctx *context.Context, event *Event) (interface{}, error) {
	log.Println("processing event")
	// Perform business logic  . . .

	return nil, nil
}
