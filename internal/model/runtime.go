package model

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// SampleEnvelope is the mutable transport object used while decoding samples.
type SampleEnvelope struct {
	Tenant  string
	Metric  string
	Labels  map[string]string
	Payload []byte
}

func (e *SampleEnvelope) Reset() {
	e.Tenant = ""
	e.Metric = ""
	clear(e.Labels)
	e.Payload = e.Payload[:0]
}

func (e *SampleEnvelope) Clone() SampleEnvelope {
	clone := SampleEnvelope{Tenant: e.Tenant, Metric: e.Metric}
	clone.Payload = append([]byte(nil), e.Payload...)
	clone.Labels = make(map[string]string, len(e.Labels))
	for key, value := range e.Labels {
		clone.Labels[key] = value
	}
	return clone
}

// EnvelopePool provides deterministic reuse without exposing sync.Pool timing.
type EnvelopePool struct {
	mu   sync.Mutex
	free *SampleEnvelope
}

func (p *EnvelopePool) Acquire() *SampleEnvelope {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.free != nil {
		e := p.free
		p.free = nil
		return e
	}
	return &SampleEnvelope{Labels: make(map[string]string)}
}

func (p *EnvelopePool) Release(e *SampleEnvelope) {
	e.Reset()
	p.mu.Lock()
	p.free = e
	p.mu.Unlock()
}

type AlertSnapshot struct {
	Values map[string]float64
	Ready  bool
}

func (s *AlertSnapshot) Clone() *AlertSnapshot {
	if s == nil {
		return nil
	}
	clone := &AlertSnapshot{Ready: s.Ready, Values: make(map[string]float64, len(s.Values))}
	for key, value := range s.Values {
		clone.Values[key] = value
	}
	return clone
}

func BuildAlertSnapshot(values map[string]*float64) (*AlertSnapshot, error) {
	built := make(map[string]float64, len(values))
	for key, value := range values {
		if value == nil {
			return nil, fmt.Errorf("alert value %q is missing", key)
		}
		built[key] = *value
	}
	return &AlertSnapshot{Values: built, Ready: true}, nil
}

type ParsedBatch struct {
	Metric string
	Frames [][]byte
}

func (b ParsedBatch) Clone() ParsedBatch {
	clone := ParsedBatch{Metric: b.Metric, Frames: make([][]byte, len(b.Frames))}
	for i := range b.Frames {
		clone.Frames[i] = append([]byte(nil), b.Frames[i]...)
	}
	return clone
}

type ReusableParser struct {
	buffer []byte
}

func (p *ReusableParser) Parse(metric string, frames ...string) ParsedBatch {
	total := 0
	for _, frame := range frames {
		total += len(frame)
	}
	if cap(p.buffer) < total {
		p.buffer = make([]byte, total)
	} else {
		p.buffer = p.buffer[:total]
	}
	offset := 0
	parsed := ParsedBatch{Metric: metric, Frames: make([][]byte, 0, len(frames))}
	for _, frame := range frames {
		copy(p.buffer[offset:], frame)
		parsed.Frames = append(parsed.Frames, p.buffer[offset:offset+len(frame)])
		offset += len(frame)
	}
	return parsed
}

type TaskUpdate struct {
	ID           string
	Version      uint64
	Status       string
	OperationKey string
}

func (u TaskUpdate) NewerThan(current TaskUpdate) bool {
	return u.Version > current.Version
}

type Validator interface {
	Validate(map[string]string) error
}

type RuleValidator struct {
	required string
}

func NewRuleValidator(enabled bool, required string) Validator {
	if !enabled {
		return nil
	}
	return &RuleValidator{required: required}
}

func (v *RuleValidator) Validate(values map[string]string) error {
	if v == nil {
		return nil
	}
	if strings.TrimSpace(values[v.required]) == "" {
		return fmt.Errorf("required setting %q is empty", v.required)
	}
	return nil
}

var (
	ErrSinkRejected  = errors.New("sink rejected sample batch")
	ErrSinkTemporary = errors.New("sink temporarily unavailable")
)

type SinkError struct {
	Kind string
	Err  error
}

func (e *SinkError) Error() string { return e.Kind + ": " + e.Err.Error() }
func (e *SinkError) Unwrap() error { return e.Err }

func NewRejectedSinkError(err error) error {
	return &SinkError{Kind: "rejected", Err: errors.Join(ErrSinkRejected, err)}
}

func NewTemporarySinkError(err error) error {
	return &SinkError{Kind: "temporary", Err: errors.Join(ErrSinkTemporary, err)}
}

func IsTemporarySinkError(err error) bool { return errors.Is(err, ErrSinkTemporary) }
