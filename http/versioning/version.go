package versioning

import (
	"time"
)

// Version represents a version change. It is pinned to a specific time.
// It contains a list of Changes. Changes are executed in-order.
type Version struct {
	Date       string
	Changes    []*Change
	Deprecated bool

	DateTime time.Time
	layout   string
}

func (v *Version) String() string {
	return v.DateTime.Format(v.layout)
}

type versions []*Version

func (vs versions) Len() int           { return len(vs) }
func (vs versions) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs versions) Less(i, j int) bool { return vs[i].DateTime.Before(vs[j].DateTime) }

// Change represents a backwards-incompatible change and the actions
// required to make it compatible.
type Change struct {
	// Description of the change made. Used for documentation.
	Description string
	// Actions are a map of object type to Action. The object type
	// is determined using the reflect pkg.
	RequestAction  Action
	ResponseAction Action
}

// Action represents an action to take on a object in order to make
// it compatible. An action takes an interface as input and returns
// an updated interface.
type Action func(map[string]interface{}) map[string]interface{}
