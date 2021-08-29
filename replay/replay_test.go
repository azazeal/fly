package replay

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/azazeal/fly/env"

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

			got, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusConflict, res.StatusCode)
			assert.Equal(t, http.StatusText(http.StatusConflict)+"\n", string(got))
			assert.Equal(t, kase.exp(), res.Header.Get("Fly-Replay"))
		})
	}
}

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

func TestInRegionHandlerForRegion(t *testing.T) {
	region, state, restoreEnv := setupInRegionHandlerTest(t)
	defer restoreEnv()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := InRegionHandler(http.HandlerFunc(execute), region, state)

	h.ServeHTTP(rec, req)
	res := rec.Result()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "executed", string(body))
}

func TestInRegionHandlerForOtherRegion(t *testing.T) {
	currentRegion, state, restoreEnv := setupInRegionHandlerTest(t)
	defer restoreEnv()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	primaryRegion := otherRegion(t, currentRegion)
	h := InRegionHandler(http.HandlerFunc(execute), primaryRegion, state)

	h.ServeHTTP(rec, req)
	res := rec.Result()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, "Conflict\n", string(body))
	assert.Equal(t, http.StatusConflict, res.StatusCode)

	exp := fmt.Sprintf("region=%s;state=%s", primaryRegion, state)
	assert.Equal(t, res.Header.Get("fly-replay"), exp)
}

func setupInRegionHandlerTest(t *testing.T) (region, state string, fn func()) {
	t.Helper()

	region = testutil.HexString(t, 3)
	fn = testutil.SetEnv(t, map[string]string{env.RegionKey: region})
	state = testutil.HexString(t, 20)

	return
}

func otherRegion(t *testing.T, region string) (other string) {
	for other = region; other == region; other = testutil.HexString(t, 3) {
		continue
	}

	return
}

func execute(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, "executed")

	return
}
