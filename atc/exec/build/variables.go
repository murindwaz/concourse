package build

import (
	"sync"

	"github.com/concourse/concourse/vars"
)

type Variables struct {
	parentScope vars.Variables
	localVars   vars.StaticVariables
	tracker     *vars.Tracker
	lock        sync.RWMutex
}

func NewVariables(tracker *vars.Tracker) *Variables {
	return &Variables{
		localVars: vars.StaticVariables{},
		tracker:   tracker,
	}
}

func (v *Variables) Get(ref vars.Reference) (interface{}, bool, error) {
	if ref.Source == "." {
		v.lock.RLock()
		val, found, err := v.localVars.Get(ref.WithoutSource())
		v.lock.RUnlock()
		if found || err != nil {
			return val, found, err
		}
	}

	if v.parentScope != nil {
		val, found, err := v.parentScope.Get(ref)
		if found || err != nil {
			return val, found, err
		}
	}

	return nil, false, nil
}

func (v *Variables) List() ([]vars.Reference, error) {
	list, err := v.parentScope.List()
	if err != nil {
		return nil, err
	}
	v.lock.RLock()
	defer v.lock.RUnlock()
	for k := range v.localVars {
		list = append(list, vars.Reference{Source: ".", Path: k})
	}
	return list, nil
}

func (v *Variables) NewScope(tracker *vars.Tracker) *Variables {
	return &Variables{
		parentScope: v,
		localVars:   vars.StaticVariables{},
		tracker:     tracker,
	}
}

func (v *Variables) SetVar(source, name string, val interface{}, redact bool) {
	v.lock.Lock()
	v.localVars[name] = val
	v.lock.Unlock()

	if redact {
		v.tracker.Track(vars.Reference{Source: source, Path: name}, val)
	}
}

func (v *Variables) RedactionEnabled() bool {
	return v.tracker.Enabled
}
