// Package domain holds the core types shared by every storage backend and
// every transport (HTTP, CLI, whatever comes later). Nothing in here knows
// about SQL, JSON files, or HTTP — that's the point: swap any of those out
// without touching this package.
package domain

import "time"

// Kind is what's being shared. It's a string, not a fixed set of Go types,
// so a new kind (a seed library? a book?) is a config change, not a code
// change — see BLUEPRINT.md "Adding a new kind of resource".
type Kind string

const (
	KindTool     Kind = "tool"
	KindHardware Kind = "hardware"
	KindSoftware Kind = "software"
)

func (k Kind) Valid() bool {
	switch k {
	case KindTool, KindHardware, KindSoftware:
		return true
	default:
		return false
	}
}

// Resource is anything a community stewards and lends out: a physical tool,
// a piece of hardware (a 3D printer, a diagnostic kit), or a software
// project (a maintained build, a licensed seat, a hosted instance).
//
// Metadata carries kind-specific detail (a tool's condition, a repo's URL,
// a hardware kit's serial number) without forcing every kind into the same
// rigid columns — this is what "parameterized, not hardcoded" (CLAUDE.md
// rule 3) looks like in a data model.
type Resource struct {
	ID          string            `json:"id"`
	Kind        Kind              `json:"kind"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

func (r Resource) Validate() error {
	if r.Name == "" {
		return ErrInvalidInput("name is required")
	}
	if !r.Kind.Valid() {
		return ErrInvalidInput("kind must be tool, hardware, or software")
	}
	return nil
}
