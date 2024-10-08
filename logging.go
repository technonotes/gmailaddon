/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
)

const (
	LevelCritical = slog.Level(12)
)

// Middleware that adds the Cloud Trace ID to the context
// This is used to correlate the structured logs with the Cloud Run
// request log.
func WithCloudTraceContext(h http.Handler) http.Handler {
	// Get the project ID from the environment if specified
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		var err error
		// Get from metadata server
		// You can avoid this dependency by using environment variables, or by connecting
		// to the metadata endpoint directly using an `http.Client`
		// See https://cloud.google.com/compute/docs/metadata/overview
		projectID, err = metadata.ProjectID()
		if err != nil {
			panic(err)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var trace string
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "trace", trace)))
	})
}

func traceFromContext(ctx context.Context) string {
	trace := ctx.Value("trace")
	if trace == nil {
		return ""
	}
	return trace.(string)
}

// Handler that outputs JSON understood by the structured log agent.
// See https://cloud.google.com/logging/docs/agent/logging/configuration#special-fields
type CloudLoggingHandler struct{ handler slog.Handler }

func NewCloudLoggingHandler() *CloudLoggingHandler {
	return &CloudLoggingHandler{handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			} else if a.Key == slog.SourceKey {
				a.Key = "logging.googleapis.com/sourceLocation"
			} else if a.Key == slog.LevelKey {
				a.Key = "severity"
				level := a.Value.Any().(slog.Level)
				if level == LevelCritical {
					a.Value = slog.StringValue("CRITICAL")
				}
			}
			return a
		},
	})}
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CloudLoggingHandler) Handle(ctx context.Context, rec slog.Record) error {
	trace := traceFromContext(ctx)
	if trace != "" {
		rec = rec.Clone()
		// Add trace ID	to the record so it is correlated with the Cloud Run request log
		// See https://cloud.google.com/trace/docs/trace-log-integration
		rec.Add("logging.googleapis.com/trace", slog.StringValue(trace))
	}
	return h.handler.Handle(ctx, rec)
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithGroup(name)}
}
