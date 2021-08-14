// Package replay implements helpers for when replaying requests.
package replay

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Header denotes fly's replay header.
const Header = "Fly-Replay"

// Set sets the Fly-Replay header of h to carry to region and optionally, state.
func Set(h http.Header, region, state string) {
	if state == "" {
		h.Set(Header, fmt.Sprintf("region=%s", region))
	} else {
		h.Set(Header, fmt.Sprintf("region=%s;state=%s", region, state))
	}
}

// In writes to w a http.StatusConflict response with a Fly-Replay header that
// carries region and optionally, state.
func In(w http.ResponseWriter, region, state string) {
	Set(w.Header(), region, state)

	const code = http.StatusConflict
	http.Error(w, http.StatusText(code), code)
}

// SourceInfo wraps the properties of a replayed request.
type SourceInfo struct {
	// Instance denotes the instance that requested the replay.
	Instance string

	// Region denotes the region the replay was initially requested in.
	Region string

	// Time denotes the time at which the replay was initially requested.
	Time time.Time

	// State denotes the user-defined state the replay carries.
	State string
}

// SourceHeader denotes fly's replay source header.
const SourceHeader = Header + "-Src"

// Source returns the replay source r carries.
//
// Source returns nil for requests that have not been replayed.
func Source(r *http.Request) (inf *SourceInfo) {
	values, ok := r.Header[Header+"-Src"]
	if !ok {
		return
	}

	inf = &SourceInfo{}

	val := values[0]
	for len(val) > 0 {
		i := strings.Index(val, ";")
		if i < 0 {
			i = len(val)
		}

		switch token := val[:i]; {
		case strings.HasPrefix(token, "instance="):
			inf.Instance = token[9:]
		case strings.HasPrefix(token, "region="):
			inf.Region = token[7:]
		case strings.HasPrefix(token, "t="):
			if μS, err := strconv.ParseInt(token[2:], 10, 64); err == nil {
				inf.Time = time.Unix(0, μS*1000)
			}
		case strings.HasPrefix(token, "state="):
			inf.State = token[6:]
		}

		if i == len(val) {
			break
		}

		val = val[i+1:]
	}

	return
}
