package replay

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/azazeal/fly/internal/testutil"
)

type inTestCase struct {
	region string
	state  string
}

func (itc *inTestCase) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	In(w, itc.region, itc.state)
}

func (itc *inTestCase) exp() string {
	tokens := []string{"region=" + itc.region}
	if itc.state != "" {
		tokens = append(tokens, "state="+itc.state)
	}
	return strings.Join(tokens, ";")
}

func TestIn(t *testing.T) {
	cases := []*inTestCase{
		0: {},
		1: {
			region: testutil.HexString(t, 3),
		},
		2: {
			state: testutil.HexString(t, 10),
		},
		3: {
			region: testutil.HexString(t, 3),
			state:  testutil.HexString(t, 20),
		},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/state", nil)

			kase.ServeHTTP(rec, req)
			res := rec.Result()

			// TODO: use io.ReadAll when support for Go 1.15 is dropped
			got, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusConflict, res.StatusCode)
			assert.Equal(t, http.StatusText(http.StatusConflict)+"\n", string(got))
			assert.Equal(t, kase.exp(), res.Header.Get("Fly-Replay"))
		})
	}
}

// Fly-Replay-Src [instance=50f17653;region=ams;t=1628960671238231;state=123]
//
// - instance: 50f17653
// - region: ams
// - time: 2021-08-14 17:04:31.238231 +0000 UTC
// - state: 123
//

func TestSource(t *testing.T) {
	buildRequest := func(add bool, hdr string) (r *http.Request) {
		r = httptest.NewRequest(http.MethodGet, "/", nil)
		if add {
			r.Header.Add("fly-replay-src", hdr)
		}

		return r
	}

	cases := []struct {
		req *http.Request
		exp *SourceInfo
	}{
		0: {
			req: buildRequest(false, ""),
		},
		1: {
			req: buildRequest(true, ""),
			exp: &SourceInfo{},
		},
		2: {
			req: buildRequest(true, "invalid values"),
			exp: &SourceInfo{},
		},
		3: {
			req: buildRequest(true, "instance=;region=;t=;state="),
			exp: &SourceInfo{},
		},
		4: {
			req: buildRequest(true, "instance=50f17653;region=ams;t=1628960671238231;state=some-state"),
			exp: &SourceInfo{
				Instance: "50f17653",
				Region:   "ams",
				Time:     time.Unix(0, int64(1628960671238231*time.Microsecond)),
				State:    "some-state",
			},
		},
		5: {
			req: buildRequest(true, "state;state=some-state;region=ams;ins;instance=50f176;53;;t=1628960671238231;"),
			exp: &SourceInfo{
				Instance: "50f176",
				Region:   "ams",
				Time:     time.Unix(0, int64(1628960671238231*time.Microsecond)),
				State:    "some-state",
			},
		},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			got := Source(kase.req)
			if kase.exp == nil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, kase.exp, got)
			}
		})
	}
}
