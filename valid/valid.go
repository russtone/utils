package valid

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/sys/unix"

	"github.com/russtone/utils/iprange"
)

//
// File
//

func FileRead() validation.Rule {
	return &fileRule{read: true}
}

func FileWrite() validation.Rule {
	return &fileRule{write: true}
}

type fileRule struct {
	read  bool
	write bool
}

func (r *fileRule) Validate(value interface{}) error {
	path, _ := value.(string)

	fi, err := os.Stat(path)
	exist := !os.IsNotExist(err)

	if r.read {
		if !exist {
			return fmt.Errorf("%q: not found", path)
		}

		if unix.Access(path, unix.R_OK) != nil {
			return fmt.Errorf("%q: permission denied", path)
		}
	} else if r.write {
		if exist && fi.IsDir() {
			return fmt.Errorf("%q: must be a file, not a directory", path)
		}

		if (exist && unix.Access(path, unix.W_OK) != nil) ||
			(!exist && unix.Access(filepath.Dir(path), unix.W_OK) != nil) {
			return fmt.Errorf("%q: permission denied", path)
		}
	}

	return nil
}

//
// Directory
//

func Directory() validation.Rule {
	return &directoryRule{}
}

type directoryRule struct{}

func (r *directoryRule) Validate(value interface{}) error {
	path, _ := value.(string)

	if fi, err := os.Stat(path); os.IsNotExist(err) {
		return err
	} else if !fi.IsDir() {
		return errors.New("must be a directory, not a file")
	}

	return nil
}

//
// OneOf
//

func OneOf(values []string, caseSensetive bool) validation.Rule {
	return &oneOfRule{values, caseSensetive}
}

type oneOfRule struct {
	values        []string
	caseSensetive bool
}

func (r *oneOfRule) Validate(value interface{}) error {
	val, _ := value.(string)

	for _, v := range r.values {
		if (r.caseSensetive && v == val) ||
			(!r.caseSensetive && strings.EqualFold(v, val)) {
			return nil
		}
	}

	return fmt.Errorf("invalid value, expected one of %s", strings.Join(r.values, ","))
}

//
// IPRange
//

func IPRange() validation.Rule {
	return &iprangeRule{}
}

type iprangeRule struct{}

func (r *iprangeRule) Validate(value interface{}) error {
	s := value.(string)
	if _, err := iprange.Parse(s); err != nil {
		return err
	}
	return nil
}

//
// Regexp
//

func Regexp() validation.Rule {
	return &regexpRule{}
}

type regexpRule struct{}

func (r *regexpRule) Validate(value interface{}) error {
	s := value.(string)
	if _, err := regexp.Compile(s); err != nil {
		return err
	}
	return nil
}
