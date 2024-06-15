package statefile

import (
	"bytes"
	"os"
	"slices"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type State struct {
	Bootstrapped []string `yaml:"bootstrapped"`
	path         string   // the original state file path
	dirty        bool     // if the state needs saving
}

func Parse(srcpath string) (*State, error) {
	var state = State{Bootstrapped: []string{}, path: srcpath}

	stateFile, err := os.Open(srcpath)

	// we don't have to worry about ErrNotExist since it just means that
	// there's no state and we can initialize it.
	if err != nil {
		if os.IsNotExist(err) {
			return &state, nil
		}
		return nil, err
	}

	defer stateFile.Close()

	err = yaml.NewDecoder(stateFile).Decode(&state)

	return &state, err
}

// check if service was not yet bootstrapped. if so, run the bootstrap function
// and update the state.
func (s *State) SetBootstrapped(service string, handler func() error) error {
	if !slices.Contains(s.Bootstrapped, service) {
		var err error
		if err = handler(); err == nil {
			s.Bootstrapped = append(s.Bootstrapped, service)
			s.Bootstrapped = lo.Uniq(s.Bootstrapped)
			s.dirty = true
		}
		return err
	}
	return nil
}

func (s *State) ClearBootstrapped(service string, handler func() error) error {
	if slices.Contains(s.Bootstrapped, service) {
		var err error
		if err = handler(); err == nil {
			s.Bootstrapped = lo.Reject(
				s.Bootstrapped,
				func(item string, _ int) bool { return item == service },
			)
			s.dirty = true
		}
		return err
	}
	return nil
}

func (s *State) Save() error {
	if !s.dirty {
		// don't save the file if we don't have to
		return nil
	}
	var buffer bytes.Buffer
	err := yaml.NewEncoder(&buffer).Encode(s)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, buffer.Bytes(), 0644)
}
