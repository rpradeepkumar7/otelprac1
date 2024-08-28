package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/attribute"
)

type LogEntry struct {
    Timestamp           string              `json:"Timestamp"`
    ObservedTimestamp   string              `json:"ObservedTimestamp"`
    TraceID             string              `json:"TraceId"`
    SpanID              string              `json:"SpanId"`
    SeverityText        string              `json:"SeverityText"`
    SeverityNumber      string              `json:"SeverityNumber"`
    Body                string              `json:"Body"`
    Resource            map[string]string   `json:"Resource"`
    InstrumentationScope map[string]string `json:"InstrumentationScope"`
    Attributes          map[string]string   `json:"Attributes"`
    EventData           map[string]string   `json:"EventData"`
    Exception           map[string]string   `json:"Exception"`
    Duration            string              `json:"Duration"`
    Status              string              `json:"Status"`
    LogLevel            string              `json:"LogLevel"`
}

func main() {
    // Set up OpenTelemetry exporter
    exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
    if err != nil {
        log.Fatal(err)
    }

    // Set up Resource with Attributes
    res, err := resource.New(
        context.Background(),
        resource.WithAttributes(
            attribute.String("service.name", "web-backend"),
            attribute.String("host.name", "web-server-1"),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Set up Trace Provider
    tracerProvider := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(res),
    )
    defer func() {
        if err := tracerProvider.Shutdown(context.Background()); err != nil {
            log.Fatal(err)
        }
    }()

    // Set the global trace provider
    otel.SetTracerProvider(tracerProvider)

    // Use the tracer (example usage)
    tracer := otel.Tracer("example-tracer")
    _, span := tracer.Start(context.Background(), "example-span")
    defer span.End()

    // Example Log Entry
    logEntry := LogEntry{
        Timestamp:         time.Now().Format(time.RFC3339),
        ObservedTimestamp: time.Now().Add(100 * time.Millisecond).Format(time.RFC3339),
        TraceID:           "abcd1234",
        SpanID:            "efgh5678",
        SeverityText:      "ERROR",
        SeverityNumber:    "17",
        Body:              "An error occurred while processing the request.",
        Resource: map[string]string{
            "service.name": "web-backend",
            "host.name":    "web-server-1",
        },
        InstrumentationScope: map[string]string{
            "Name":    "GoLogger",
            "Version": "1.0.0",
        },
        Attributes: map[string]string{
            "http.method":      "GET",
            "http.status_code": "500",
            "http.url":         "http://example.com",
            "db.operation":     "SELECT",
        },
        EventData: map[string]string{
            "event.name": "request_error",
            "event.type": "error",
        },
        Exception: map[string]string{
            "exception.message":  "Database connection failed",
            "exception.type":     "DatabaseError",
            "exception.stacktrace": "at com.example.Database.connect(Database.java:42)\n...more stack trace...",
        },
        Duration: "100ms",
        Status:   "failed",
        LogLevel: "error",
    }

    // Convert log entry to JSON and print it
    logEntryJSON, err := json.MarshalIndent(logEntry, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Log Entry in JSON format:")
    log.Println(string(logEntryJSON))

    // Your application logic here
    log.Println("OpenTelemetry is set up and running!")
}
