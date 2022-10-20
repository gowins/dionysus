package registry

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/gowins/dionysus/recovery"
)

var registries sync.Map
var Default Registry

func Register(name string, r Registry) error {
	if name == "" || r == nil {
		return fmt.Errorf("[name] or [r] is nil")
	}

	_, ok := registries.LoadOrStore(name, r)
	if ok {
		return fmt.Errorf("registry %s is exists", name)
	}
	return nil
}

func Init(rawUrl string) (gErr error) {
	defer recovery.CheckErr(&gErr)

	url1, err := url.Parse(rawUrl)
	if err != nil {
		return errors.Wrapf(err, "url %s parse error", rawUrl)
	}

	scheme := url1.Scheme
	if Get(scheme) == nil {
		return fmt.Errorf("registry [%s] not exists", scheme)
	}

	params := url1.Query()
	var opts []Option
	if val := params.Get("secure"); val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return errors.Wrapf(err, "secure %s ParseBool error", val)
		}
		opts = append(opts, Secure(b))
	}

	if val := params.Get("timeout"); val != "" {
		dur, err := time.ParseDuration(val)
		if err != nil {
			return errors.Wrapf(err, "timeout %s ParseDuration error", val)
		}
		opts = append(opts, Timeout(dur))
	}

	if val := params.Get("ttl"); val != "" {
		dur, err := time.ParseDuration(val)
		if err != nil {
			return errors.Wrapf(err, "ttl %s ParseDuration error", val)
		}
		opts = append(opts, TTL(dur))
	}

	if strings.TrimSpace(url1.Host) == "" {
		return errors.New("host is null")
	}

	addrs := strings.Split(url1.Host, ",")
	if len(addrs) == 0 {
		return fmt.Errorf("%s host should be nil", rawUrl)
	}

	opts = append(opts, Addrs(addrs...))

	Default = Get(scheme)
	if err := Default.Init(opts...); err != nil {
		return errors.Wrapf(err, "[%s] init error", scheme)
	}

	return nil
}

func Get(name string) Registry {
	val, ok := registries.Load(name)
	if !ok {
		return nil
	}

	return val.(Registry)
}
