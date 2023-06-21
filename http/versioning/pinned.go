package versioning

import (
	"errors"
	"net/http"
	"sort"
	"time"
)

var (
	// ErrInvalidVersion means the version supplied is not included
	// in the versions supplied to the VersionManager or it is malformed.
	ErrInvalidVersion = errors.New("invalid version")

	// ErrNoVersionSupplied means no version was supplied.
	ErrNoVersionSupplied = errors.New("no version supplied")

	// ErrVersionDeprecated means the version is deprecated.
	ErrVersionDeprecated = errors.New("version is deprecated")
)

const (
	defaultLayout = "2006-01-02"
	defaultHeader = "api_version"
	defaultQuery  = "v"
)

// VersionManager represents a list of versions.
type VersionManager struct {
	Layout string
	Header string
	Query  string
}

func (vm *VersionManager) ParseVersion(v *Version) error {
	var err error
	v.layout = vm.layout()
	v.DateTime, err = time.Parse(vm.layout(), v.Date)
	if err != nil {
		return err
	}

	return nil
}

// Latest returns the most current active version.
func (vm *VersionManager) Latest(vs []*Version) *Version {
	if len(vs) == 0 {
		return nil
	}
	return vs[0]
}

// It inspects the query parameters and request headers. Whichever
// is most recent is the version to use.
func (vm *VersionManager) Parse(r *http.Request, vs []*Version) (*Version, error) {
	for _, v := range vs {
		if err := vm.ParseVersion(v); err != nil {
			return nil, err
		}
	}
	sort.Sort(sort.Reverse(versions(vs)))

	h := r.Header.Get(vm.header())
	q := r.URL.Query().Get(vm.query())

	if h == "" && q == "" {
		return nil, ErrNoVersionSupplied
	}

	hDate, qDate := time.Time{}, time.Time{}

	var err error
	if h != "" {
		hDate, err = time.Parse(vm.layout(), h)
		if err != nil {
			return nil, ErrInvalidVersion
		}
	}
	if q != "" {
		qDate, err = time.Parse(vm.layout(), q)
		if err != nil {
			return nil, ErrInvalidVersion
		}
	}

	t := hDate
	if hDate.Before(qDate) {
		t = qDate
	}

	v, err := vm.getVersionByTime(t, vs)
	if err != nil {
		return nil, err
	}
	if v.Deprecated {
		return v, ErrVersionDeprecated
	}
	return v, nil
}

func (vm *VersionManager) getVersionByTime(t time.Time, versions []*Version) (*Version, error) {
	for _, v := range versions {
		if t.Equal(v.DateTime) {
			return v, nil
		}
	}

	return nil, ErrInvalidVersion
}

// ApplyRequest processes a object by applying all changes between the
// latest version and the version requested. The altered object is returned.
//
// Concretely, if the supplied version is two versions behind the latest, the changes
// in those two versions are applied sequentially to the object. This essentially
// "undoes" the changes made to the API so that the object is structured according to
// the specified version.
func (vm *VersionManager) ApplyRequest(obj map[string]interface{}, version *Version, versions []*Version) (map[string]interface{}, error) {

	requestedVersionDate := version.DateTime
	for _, v := range versions {
		// If the requested version is >= to the version, do not apply.
		if requestedVersionDate.After(v.DateTime) || requestedVersionDate.Equal(v.DateTime) {
			break
		}
		// Iterate through each change and execute
		// actions as appropriate.
		for _, c := range v.Changes {
			// If there is an action for this obj type
			// execute the action.
			if nil == c.RequestAction {
				return nil, nil
			}
			obj = c.RequestAction(obj)
		}

	}

	return obj, nil
}

func (vm *VersionManager) ApplyResponse(obj map[string]interface{}, version *Version, versions []*Version) (map[string]interface{}, error) {
	requestedVersionDate := version.DateTime
	for _, v := range versions {

		// If the requested version is >= to the version, do not apply.

		if requestedVersionDate.After(v.DateTime) || requestedVersionDate.Equal(v.DateTime) {
			break
		}

		// Iterate through each change and execute
		// actions as appropriate.
		for _, c := range v.Changes {
			// If there is an action for this obj type
			// execute the action.
			if nil == c.ResponseAction {
				return nil, nil
			}
			obj = c.ResponseAction(obj)
		}
	}

	return obj, nil
}

// Versions returns a list of all versions as strings.
func (vm *VersionManager) Versions(vs []*Version) []string {
	versions := make([]string, len(vs))
	for i := range versions {
		versions[i] = vs[i].String()
	}

	return versions
}

func (vm *VersionManager) layout() string {
	if vm.Layout != "" {
		return vm.Layout
	}

	return defaultLayout
}

func (vm *VersionManager) header() string {
	if vm.Header != "" {
		return vm.Header
	}

	return defaultHeader
}

func (vm *VersionManager) query() string {
	if vm.Query != "" {
		return vm.Query
	}

	return defaultQuery
}
