package envelope

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"
)

var hostName string

func init() {
	h, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	} else {
		hostName = h
	}
}

type Envelope[T proto.Message] struct {
	Payload   T
	Origin    string
	Timestamp int64
}

// Pack creates a new envelope with hostname and timestamp set.
func Pack[T proto.Message](payload T) *Envelope[T] {
	return &Envelope[T]{
		Payload:   payload,
		Origin:    hostName,
		Timestamp: time.Now().UnixMilli(),
	}
}

// ToDict serializes payload to raw bytes (for Redis or MQ).
func (e *Envelope[T]) ToMap() (map[string]any, error) {
	raw, err := proto.Marshal(e.Payload)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"payload":   raw,
		"origin":    e.Origin,
		"timestamp": e.Timestamp,
	}, nil
}

// ToJSONSafe returns a map with base64-encoded payload for JSON use.
func (e *Envelope[T]) ToJSONSafe() (map[string]any, error) {
	raw, err := proto.Marshal(e.Payload)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"payload":   base64.StdEncoding.EncodeToString(raw),
		"origin":    e.Origin,
		"timestamp": e.Timestamp,
	}, nil
}

// FromJSON parses a generic map into an Envelope.
func FromJSON[T proto.Message](data map[string]any, target T) (*Envelope[T], error) {

	// payload: []byte or string
	var payload []byte
	switch v := data["payload"].(type) {
	case []byte:
		payload = v
	case string:
		payload = []byte(v)
	default:
		return nil, fmt.Errorf("payload type error: %T", v)
	}
	if err := proto.Unmarshal(payload, target); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	// origin: string
	var origin string
	if v, ok := data["origin"].(string); ok {
		origin = v
	} else {
		return nil, fmt.Errorf("origin type error: %T", data["origin"])
	}
	// timestamp: float64, int64, int, or string
	var timestamp int64
	switch v := data["timestamp"].(type) {
	case int64:
		timestamp = v
	case int:
		timestamp = int64(v)
	case string:
		ts, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp string: %w", err)
		}
		timestamp = ts
	case float64:
		timestamp = int64(v)
	default:
		return nil, fmt.Errorf("timestamp type error: %T", v)
	}

	return &Envelope[T]{
		Payload:   target,
		Origin:    origin,
		Timestamp: (timestamp),
	}, nil
}
