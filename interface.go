package hcsshim

import (
	"encoding/json"
	"io"
	"time"
)

// ProcessConfig is used as both the input of Container.CreateProcess
// and to convert the parameters to JSON for passing onto the HCS
type ProcessConfig struct {
	ApplicationName   string            `json:",omitempty"`
	CommandLine       string            `json:",omitempty"`
	CommandArgs       []string          `json:",omitempty"` // Used by Linux Containers on Windows
	User              string            `json:",omitempty"`
	WorkingDirectory  string            `json:",omitempty"`
	Environment       map[string]string `json:",omitempty"`
	EmulateConsole    bool              `json:",omitempty"`
	CreateStdInPipe   bool              `json:",omitempty"`
	CreateStdOutPipe  bool              `json:",omitempty"`
	CreateStdErrPipe  bool              `json:",omitempty"`
	ConsoleSize       [2]uint           `json:",omitempty"`
	CreateInUtilityVm bool              `json:",omitempty"` // Used by Linux Containers on Windows
	OCISpecification  *json.RawMessage  `json:",omitempty"` // Used by Linux Containers on Windows
}

type ComputeSystemQuery struct {
	IDs    []string `json:"Ids,omitempty"`
	Types  []string `json:",omitempty"`
	Names  []string `json:",omitempty"`
	Owners []string `json:",omitempty"`
}

// Container represents a created (but not necessarily running) container.
type Container interface {
	// Start synchronously starts the container.
	Start() error

	// Shutdown requests a container shutdown, but it may not actually be shutdown until Wait() succeeds.
	Shutdown() error

	// Terminate requests a container terminate, but it may not actually be terminated until Wait() succeeds.
	Terminate() error

	// Waits synchronously waits for the container to shutdown or terminate.
	Wait() error

	// WaitTimeout synchronously waits for the container to terminate or the duration to elapse. It
	// returns false if timeout occurs.
	WaitTimeout(time.Duration) error

	// Pause pauses the execution of a container.
	Pause() error

	// Resume resumes the execution of a container.
	Resume() error

	// HasPendingUpdates returns true if the container has updates pending to install.
	HasPendingUpdates() (bool, error)

	// Statistics returns statistics for a container.
	Statistics() (Statistics, error)

	// ProcessList returns details for the processes in a container.
	ProcessList() ([]ProcessListItem, error)

	// MappedVirtualDisks returns virtual disks mapped to a utility VM, indexed by controller
	MappedVirtualDisks() (map[int]MappedVirtualDiskController, error)

	// CreateProcess launches a new process within the container.
	CreateProcess(c *ProcessConfig) (Process, error)

	// OpenProcess gets an interface to an existing process within the container.
	OpenProcess(pid int) (Process, error)

	// Close cleans up any state associated with the container but does not terminate or wait for it.
	Close() error

	// Modify the System
	Modify(config *ResourceModificationRequestResponse) error
}

// Process represents a running or exited process.
type Process interface {
	// Pid returns the process ID of the process within the container.
	Pid() int

	// Kill signals the process to terminate but does not wait for it to finish terminating.
	Kill() error

	// Wait waits for the process to exit.
	Wait() error

	// WaitTimeout waits for the process to exit or the duration to elapse. It returns
	// false if timeout occurs.
	WaitTimeout(time.Duration) error

	// ExitCode returns the exit code of the process. The process must have
	// already terminated.
	ExitCode() (int, error)

	// ResizeConsole resizes the console of the process.
	ResizeConsole(width, height uint16) error

	// Stdio returns the stdin, stdout, and stderr pipes, respectively. Closing
	// these pipes does not close the underlying pipes; it should be possible to
	// call this multiple times to get multiple interfaces.
	Stdio() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error)

	// CloseStdin closes the write side of the stdin pipe so that the process is
	// notified on the read side that there is no more data in stdin.
	CloseStdin() error

	// Close cleans up any state associated with the process but does not kill
	// or wait on it.
	Close() error
}
