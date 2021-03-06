package runnerpool

import (
	"context"
	"io"
	"net/http"

	"github.com/fnproject/fn/api/common"
	"github.com/fnproject/fn/api/models"
)

// Placer implements a placement strategy for calls that are load-balanced
// across runners in a pool
type Placer interface {
	PlaceCall(rp RunnerPool, ctx context.Context, call RunnerCall) error
}

// RunnerPool is the abstraction for getting an ordered list of runners to try for a call
type RunnerPool interface {
	// returns an error for unrecoverable errors that should not be retried
	Runners(call RunnerCall) ([]Runner, error)
	Shutdown(ctx context.Context) error
}

// PKIData encapsulates TLS certificate data
type PKIData struct {
	Ca   string
	Key  string
	Cert string
}

// MTLSRunnerFactory represents a factory method for constructing runners using mTLS
type MTLSRunnerFactory func(addr, certCommonName string, pki *PKIData) (Runner, error)

// RunnerStatus is general information on Runner health as returned by Runner::Status() call
type RunnerStatus struct {
	ActiveRequestCount int32           // Number of active running requests on Runner
	StatusFailed       bool            // True if Status execution failed
	StatusId           string          // Call ID for Status
	Details            string          // General/Debug Log information
	ErrorCode          int32           // If StatusFailed, then error code is set
	ErrorStr           string          // Error details if StatusFailed and ErrorCode is set
	CreatedAt          common.DateTime // Status creation date at Runner
	StartedAt          common.DateTime // Status execution date at Runner
	CompletedAt        common.DateTime // Status completion date at Runner
}

// Runner is the interface to invoke the execution of a function call on a specific runner
type Runner interface {
	TryExec(ctx context.Context, call RunnerCall) (bool, error)
	Status(ctx context.Context) (*RunnerStatus, error)
	Close(ctx context.Context) error
	Address() string
}

// RunnerCall provides access to the necessary details of request in order for it to be
// processed by a RunnerPool
type RunnerCall interface {
	SlotHashId() string
	Extensions() map[string]string
	RequestBody() io.ReadCloser
	ResponseWriter() http.ResponseWriter
	StdErr() io.ReadWriteCloser
	Model() *models.Call
}
