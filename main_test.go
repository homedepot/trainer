package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homedepot/trainer/api"
	"github.com/homedepot/trainer/cli"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/handler"
	"github.com/homedepot/trainer/router"
	"github.com/homedepot/trainer/structs/state"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/juju/loggo"
	"github.com/stretchr/testify/assert"
)

var Username = "blahdeblah"
var Password = "the penguin is walking on the hamburger"
var Config *config.Config
var Handler *handler.Handler

func TestMain(m *testing.M) {
	// The tempdir is created so MongoDB has a location to store its files.
	// Contents are wiped once the server stops
	err := os.Setenv("APIAUTHUSERNAME", Username)
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	err = os.Setenv("APIAUTHPASSWORD", Password)
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	err = os.Setenv("TESTMODE", "true")
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	err = os.Setenv("LOGLEVEL", "TRACE")
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	err = os.Setenv("CONFIGFILE", "data/config.yml")
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	o = *cli.NewOptions()
	o.Parse([]string{})

	Config = LoadConfig(o)
	loglevel, berr := loggo.ParseLevel(o.LogLevel)
	if berr == false {
		log.Printf("Couldn't parse loglevel %s", o.LogLevel)
		return
	}
	loggo.GetLogger("default").SetLogLevel(loglevel)

	// Run the test suite
	m.Run()
}

func FailWithMessage(t *testing.T, message string) {
	log.Printf(message)
	t.Fail()
}

func InitHandler(t *testing.T) {
	Handler = &handler.Handler{}
	Handler.Start()
}

func StopHandler(t *testing.T) {
	Handler.Stop()
	Handler.Reset()
}

func MakeRequest(opt cli.Options, method string, url string, payload string, t *testing.T, code int) (*gin.Context, *httptest.ResponseRecorder) {
	logger := loggo.GetLogger("default")
	r := router.Router{}
	res := httptest.NewRecorder()
	c, e := gin.CreateTestContext(res)
	r.CreateRouter(opt, Config, logger, e, Handler)
	req := httptest.NewRequest(method, url, strings.NewReader(payload))
	c.Request = req
	c.Request.SetBasicAuth(Username, Password)
	c.Request.Host = "127.0.0.1"
	r.Router.ServeHTTP(res, c.Request)
	if res.Code != code {
		log.Printf("%d != %d", res.Code, code)
		log.Printf("location: %s", res.Header().Get("Location"))
		t.Errorf("wrong code: %s", strconv.Itoa(res.Code))
	}
	return c, res
}

func MakeRequestHttp(opt cli.Options, method string, url string, payload string, t *testing.T, code int) (*gin.Context, *httptest.ResponseRecorder) {
	logger := loggo.GetLogger("default")
	r := router.Router{}
	res := httptest.NewRecorder()
	c, e := gin.CreateTestContext(res)
	r.CreateRouter(opt, Config, logger, e, Handler)
	req := httptest.NewRequest(method, url, strings.NewReader(payload))
	c.Request = req
	c.Request.SetBasicAuth(Username, Password)
	c.Request.Header.Add("X-Forwarded-Proto", "http")
	c.Request.Host = "127.0.0.1"
	e.ServeHTTP(res, c.Request)
	if res.Code != code {
		log.Printf("%d != %d", res.Code, code)
		log.Printf("location: %s", res.Header().Get("Location"))
		t.Errorf("wrong code: %s", strconv.Itoa(res.Code))
	}
	return c, res
}

func MakeRequestHttps(opt cli.Options, method string, url string, payload string, t *testing.T, code int) (*gin.Context, *httptest.ResponseRecorder) {
	logger := loggo.GetLogger("default")
	r := router.Router{}
	res := httptest.NewRecorder()
	c, e := gin.CreateTestContext(res)
	r.CreateRouter(opt, Config, logger, e, Handler)
	req := httptest.NewRequest(method, url, strings.NewReader(payload))
	c.Request = req
	c.Request.SetBasicAuth(Username, Password)
	c.Request.Header.Add("X-Forwarded-Proto", "https")
	c.Request.Host = "127.0.0.1"
	e.ServeHTTP(res, c.Request)
	if res.Code != code {
		log.Printf("%d != %d", res.Code, code)
		log.Printf("location: %s", res.Header().Get("Location"))
		t.Errorf("wrong code: %s", strconv.Itoa(res.Code))
	}
	return c, res
}

func workingDir(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to find working directory: %v", err)
	}
	return wd
}

func TestConfig(t *testing.T) {

	logger := loggo.GetLogger("default")
	logger.SetLogLevel(loggo.DEBUG)
	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))
	assert.Equal(t, 27, len(Config.Plans), "Wrong number of plans")
	p := Config.Plans[0]
	assert.Equal(t, "basic_test", p.Name, "Wrong name")
	assert.Equal(t, 12, len(p.Txn), "Wrong number of transactions")
	//t0 := p.Txn[0]
	t0 := p.Txn[0]
	t1 := p.Txn[1]

	//assert.Equal(t, "init_action", t0.Name, "Wrong transaction name")
	assert.Equal(t, "transaction_1", t0.Name, "Wrong transaction name")
	assert.Equal(t, "transaction_2", t1.Name, "Wrong transaction name")
	assert.Equal(t, "/api/v1/url1", t0.URL, "Wrong URL")
	assert.Equal(t, "/api/v1/url2", t1.URL, "Wrong URL")
	oe1 := t0.OnExpected
	oe2 := t1.OnExpected
	ou1 := t0.OnUnexpected
	ou2 := t1.OnUnexpected

	assert.Equal(t, "url1.yml", oe1.Response)
	assert.Equal(t, "url2.yml", oe2.Response)
	assert.Equal(t, "url1_unexpected.yml", ou1.Response)
	assert.Equal(t, "url2_unexpected.yml", ou2.Response)

	assert.Equal(t, "200", oe1.ResponseCode)
	assert.Equal(t, "401", ou1.ResponseCode)
	assert.Equal(t, "200", oe2.ResponseCode)
	assert.Equal(t, "401", ou2.ResponseCode)

}

func TestLaunch(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/basic_test", "", t, 200)

}

func TestGetFirstTransaction(t *testing.T) {

	defer func(d string) { os.Chdir(d) }(workingDir(t))

	txn, err := Config.GetFirstTransaction("basic_test")
	if err != nil {
		FailWithMessage(t, err.Error())
		return
	}
	assert.Equal(t, "check_testvar", txn.Name, "Incorrect first transaction name")

}

/*func TestReset(t *testing.T) {
	core.S.RunnerKillSwitch = true
	time.Sleep(1 * time.Second)
	core.S.RunnerKillSwitch = false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "guid": "000000", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				FailWithMessage(t, "Couldn't read body: "+err.Error())
			}
			var r interface{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				FailWithMessage(t, "couldn't marshal json: "+err.Error())
			}
			m := r.(map[string]interface{})
			s, ok := m["Michael"]
			if !ok {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if s != "Scott" {
				FailWithMessage(t, "JSON import failed")
				return
			}
			log.Printf("Callback succeeded!")
		} else {
			FailWithMessage(t, "Got incorrect URL for callback")
			return
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL

	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/basic_test", "", t, 200)

	time.Sleep(2 * time.Second)

	_, _ = MakeRequest(o, "POST", "/capi/v1/reset", "", t, 200)

	time.Sleep(2 * time.Second)

	assert.Equal(t, "transaction_1", core.S.Transaction.Name, "Transaction didn't reset")

	v := core.S.Plan.Variables
	if v["guid"] != "000000" {
		log.Printf("variables: %+v", v)
		FailWithMessage(t, "GUID incorrect")
		return
	}
	log.Printf("variables: %+v", v)
}*/

func TestFallthrough(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/314159265358979323", "", t, 200)
}

func TestRunPlanCorrectCompleteSuccess(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "guid": "GUID", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				FailWithMessage(t, "Couldn't read body: "+err.Error())
			}
			var re interface{}
			err = json.Unmarshal(body, &re)
			if err != nil {
				FailWithMessage(t, "couldn't marshal json: "+err.Error())
			}
			m := re.(map[string]interface{})
			s, ok := m["Michael"]
			if !ok {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if s != "Scott" {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if r.Header.Get("Authorization") != "Basic bWlzdHk6Z290dGFwQHNzZW1hbGw=" {
				FailWithMessage(t, "Incorrect Authorization header: "+r.Header.Get("Authorization"))
			}
			log.Printf("Callback succeeded!")
		} else if r.URL.Path == "/GUID" {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(finishedsuccessfulresultoneci))
			log.Printf("Callback succeeded!")
		} else {
			log.Printf("Got incorrect URL %s for callback", r.URL.Path)
			t.Errorf("Got incorrect URL %s for callback", r.URL.Path)
			return
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL

	//core.S.RunnerKillSwitch = true
	//time.Sleep(1 * time.Second)
	//core.S.RunnerKillSwitch = false

	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, resp := MakeRequest(o, "POST", "/capi/v1/launch/basic_test", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	_, _ = MakeRequest(o, "POST", "/api/v1/url1", `{"amount": "too much", "I": "love tacos"}`, t, 200)
	_, _ = MakeRequest(o, "POST", "/api/v1/url2", "oompa loompa doopity doo", t, 200)

	p := handler.GetPlan()
	// let the rest of the stuff run
	time.Sleep(10 * time.Second)

	st := p.State
	assert.Equal(t, "empty", p.State.Transaction, "wth?")
	assert.Equal(t, "completed", st.States[0].Status, "Incorrect status %s", st.States[0].Status)
	assert.Equal(t, "completed", st.States[1].Status, "Incorrect status %s", st.States[1].Status)
	assert.Equal(t, "completed", st.States[2].Status, "Incorrect status %s", st.States[2].Status)
	assert.Equal(t, "expected", st.States[3].Status, "Incorrect status %s", st.States[3].Status)
	assert.Equal(t, "expected", st.States[4].Status, "Incorrect status %s", st.States[4].Status)
	assert.Equal(t, "completed", st.States[5].Status, "Incorrect status %s", st.States[5].Status)
	assert.Equal(t, "completed", st.States[6].Status, "Incorrect status %s", st.States[6].Status)
	assert.Equal(t, "completed", st.States[7].Status, "Incorrect status %s", st.States[7].Status)
	assert.Equal(t, "completed", st.States[8].Status, "Incorrect status %s", st.States[8].Status)
	assert.Equal(t, "stopped", st.States[9].Status, "Incorrect status %s", st.States[9].Status)
	assert.Equal(t, "check_testvar", st.States[0].TxnName, "Incorrect Transaction %s", st.States[0].TxnName)
	assert.Equal(t, "set_testvar", st.States[1].TxnName, "Incorrect Transaction %s", st.States[1].TxnName)
	assert.Equal(t, "init_action", st.States[2].TxnName, "Incorrect Transaction %s", st.States[2].TxnName)
	assert.Equal(t, "transaction_1", st.States[3].TxnName, "Incorrect Transaction %s", st.States[3].TxnName)
	assert.Equal(t, "transaction_2", st.States[4].TxnName, "Incorrect Transaction %s", st.States[4].TxnName)
	assert.Equal(t, "test_finish", st.States[5].TxnName, "Incorrect Transaction %s", st.States[5].TxnName)
	assert.Equal(t, "check_complete", st.States[6].TxnName, "Incorrect Transaction %s", st.States[6].TxnName)
	assert.Equal(t, "check_success", st.States[7].TxnName, "Incorrect Transaction %s", st.States[7].TxnName)
	assert.Equal(t, "txn_success", st.States[8].TxnName, "Incorrect Transaction %s", st.States[8].TxnName)
	assert.Equal(t, "empty", p.State.States[9].TxnName, "Incorrect Transaction %s", p.State.States[9].TxnName)

	_, resp = MakeRequest(o, "POST", "/capi/v1/status", "", t, 200)

	body, _ := ioutil.ReadAll(resp.Body)
	sst := state.State{}
	err := json.Unmarshal(body, &sst)
	if err != nil {
		FailWithMessage(t, err.Error())
		return
	}
	assert.Equal(t, "completed", sst.States[2].Status, "Wrong state received")
	assert.Equal(t, "expected", sst.States[3].Status, "Wrong state received")
	assert.Equal(t, "expected", sst.States[4].Status, "Wrong state received")
	assert.Equal(t, "init_action", sst.States[2].TxnName, "Wrong state received")
	assert.Equal(t, "transaction_1", sst.States[3].TxnName, "Wrong state received")
	assert.Equal(t, "transaction_2", sst.States[4].TxnName, "Wrong state received")
	assert.Equal(t, "test_finish", sst.States[5].TxnName, "Wrong state received")
	assert.Equal(t, "check_complete", sst.States[6].TxnName, "Wrong state received")
	assert.Equal(t, "check_success", sst.States[7].TxnName, "Wrong state received")
	assert.Equal(t, "txn_success", sst.States[8].TxnName, "Wrong state received")
	assert.Equal(t, "empty", sst.States[9].TxnName, "Wrong state received")
	assert.Equal(t, true, st.Variables["success"])
	assert.Equal(t, false, st.Variables["failure"])
	log.Printf("states: %+v", st.States)
}

func TestRunPlanCorrectCompleteFailure(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "guid": "GUID", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				FailWithMessage(t, "Couldn't read body: "+err.Error())
			}
			var r interface{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				FailWithMessage(t, "couldn't marshal json: "+err.Error())
			}
			m := r.(map[string]interface{})
			s, ok := m["Michael"]
			if !ok {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if s != "Scott" {
				FailWithMessage(t, "JSON import failed")
				return
			}
			log.Printf("Callback succeeded!")
		} else if r.URL.Path == "/GUID" {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(finishedunsuccessfulresultoneci))
			log.Printf("Callback succeeded!")
		} else {
			FailWithMessage(t, "Got incorrect URL for callback")
			return
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL

	//core.S.RunnerKillSwitch = true
	//time.Sleep(1 * time.Second)
	//core.S.RunnerKillSwitch = false

	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, resp := MakeRequest(o, "POST", "/capi/v1/launch/basic_test", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	p := handler.GetPlan()

	_, _ = MakeRequest(o, "POST", "/api/v1/url1", `{"amount": "too much", "I": "love tacos"}`, t, 200)
	_, _ = MakeRequest(o, "POST", "/api/v1/url2", "oompa loompa doopity doo", t, 200)

	// let the rest of the stuff run
	time.Sleep(10 * time.Second)

	st := p.State
	assert.Equal(t, "empty", p.State.Transaction, "wth?")
	assert.Equal(t, "completed", st.States[0].Status, "Incorrect status %s", st.States[0].Status)
	assert.Equal(t, "completed", st.States[1].Status, "Incorrect status %s", st.States[1].Status)
	assert.Equal(t, "completed", st.States[2].Status, "Incorrect status %s", st.States[2].Status)
	assert.Equal(t, "expected", st.States[3].Status, "Incorrect status %s", st.States[3].Status)
	assert.Equal(t, "expected", st.States[4].Status, "Incorrect status %s", st.States[4].Status)
	assert.Equal(t, "completed", st.States[5].Status, "Incorrect status %s", st.States[5].Status)
	assert.Equal(t, "completed", st.States[6].Status, "Incorrect status %s", st.States[6].Status)
	assert.Equal(t, "completed", st.States[7].Status, "Incorrect status %s", st.States[7].Status)
	assert.Equal(t, "completed", st.States[8].Status, "Incorrect status %s", st.States[8].Status)
	assert.Equal(t, "stopped", st.States[9].Status, "Incorrect status %s", st.States[9].Status)
	assert.Equal(t, "check_testvar", st.States[0].TxnName, "Incorrect Transaction %s", st.States[0].TxnName)
	assert.Equal(t, "set_testvar", st.States[1].TxnName, "Incorrect Transaction %s", st.States[1].TxnName)
	assert.Equal(t, "init_action", st.States[2].TxnName, "Incorrect Transaction %s", st.States[2].TxnName)
	assert.Equal(t, "transaction_1", st.States[3].TxnName, "Incorrect Transaction %s", st.States[3].TxnName)
	assert.Equal(t, "transaction_2", st.States[4].TxnName, "Incorrect Transaction %s", st.States[4].TxnName)
	assert.Equal(t, "test_finish", st.States[5].TxnName, "Incorrect Transaction %s", st.States[5].TxnName)
	assert.Equal(t, "check_complete", st.States[6].TxnName, "Incorrect Transaction %s", st.States[6].TxnName)
	assert.Equal(t, "check_success", st.States[7].TxnName, "Incorrect Transaction %s", st.States[7].TxnName)
	assert.Equal(t, "txn_failure", st.States[8].TxnName, "Incorrect Transaction %s", st.States[8].TxnName)
	assert.Equal(t, "empty", p.State.States[9].TxnName, "Incorrect Transaction %s", p.State.States[9].TxnName)

	_, resp = MakeRequest(o, "POST", "/capi/v1/status", "", t, 200)

	body, _ := ioutil.ReadAll(resp.Body)
	sst := state.State{}
	err := json.Unmarshal(body, &sst)
	if err != nil {
		FailWithMessage(t, err.Error())
		return
	}
	assert.Equal(t, "completed", sst.States[2].Status, "Wrong state received")
	assert.Equal(t, "expected", sst.States[3].Status, "Wrong state received")
	assert.Equal(t, "expected", sst.States[4].Status, "Wrong state received")
	assert.Equal(t, "init_action", sst.States[2].TxnName, "Wrong state received")
	assert.Equal(t, "transaction_1", sst.States[3].TxnName, "Wrong state received")
	assert.Equal(t, "transaction_2", sst.States[4].TxnName, "Wrong state received")
	assert.Equal(t, "test_finish", sst.States[5].TxnName, "Wrong state received")
	assert.Equal(t, "check_complete", sst.States[6].TxnName, "Wrong state received")
	assert.Equal(t, "check_success", sst.States[7].TxnName, "Wrong state received")
	assert.Equal(t, "txn_failure", sst.States[8].TxnName, "Wrong state received")
	assert.Equal(t, "empty", sst.States[9].TxnName, "Wrong state received")
	assert.Equal(t, false, st.Variables["success"])
	assert.Equal(t, true, st.Variables["failure"])
	log.Printf("states: %+v", st.States)
}

func TestRunPlanCorrectIncomplete(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "guid": "GUID", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				FailWithMessage(t, "Couldn't read body: "+err.Error())
			}
			var r interface{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				FailWithMessage(t, "couldn't marshal json: "+err.Error())
			}
			m := r.(map[string]interface{})
			s, ok := m["Michael"]
			if !ok {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if s != "Scott" {
				FailWithMessage(t, "JSON import failed")
				return
			}
			log.Printf("Callback succeeded!")
		} else if r.URL.Path == "/GUID" {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(unsuccessfulresultoneci))
			log.Printf("Callback succeeded!")
		} else {
			FailWithMessage(t, "Got incorrect URL for callback")
			return
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	//core.S.RunnerKillSwitch = true
	//time.Sleep(1 * time.Second)
	//core.S.RunnerKillSwitch = false

	defer func(d string) { os.Chdir(d) }(workingDir(t))
	LoadConfig(o)

	_, resp := MakeRequest(o, "POST", "/capi/v1/launch/basic_test", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	p := handler.GetPlan()

	_, _ = MakeRequest(o, "POST", "/api/v1/url1", `{"amount": "too much", "I": "love tacos"}`, t, 200)
	_, _ = MakeRequest(o, "POST", "/api/v1/url2", "oompa loompa doopity doo", t, 200)

	// let the rest of the stuff run
	time.Sleep(10 * time.Second)

	//assert.Equal(t, "check_complete", State.Transaction.Name, "wth?")
	assert.Equal(t, "completed", p.State.States[2].Status, "Incorrect status %s", p.State.States[2].Status)
	assert.Equal(t, "expected", p.State.States[3].Status, "Incorrect status %s", p.State.States[3].Status)
	assert.Equal(t, "expected", p.State.States[4].Status, "Incorrect status %s", p.State.States[4].Status)
	assert.Equal(t, "init_action", p.State.States[2].TxnName, "Incorrect Transaction %s", p.State.States[2].TxnName)
	assert.Equal(t, "transaction_1", p.State.States[3].TxnName, "Incorrect Transaction %s", p.State.States[3].TxnName)
	assert.Equal(t, "transaction_2", p.State.States[4].TxnName, "Incorrect Transaction %s", p.State.States[4].TxnName)

	_, resp = MakeRequest(o, "POST", "/capi/v1/status", "", t, 200)
	body, _ := ioutil.ReadAll(resp.Body)
	st := state.State{}
	err := json.Unmarshal(body, &st)
	if err != nil {
		FailWithMessage(t, err.Error())
		return
	}
	// This proves it looped.  No need to prove anything further at the moment, though we may want to be more
	// detailed at some point.
	assert.Equal(t, "completed", st.States[2].Status, "Wrong state received")
	assert.Equal(t, "expected", st.States[3].Status, "Wrong state received")
	assert.Equal(t, "expected", st.States[4].Status, "Wrong state received")
	assert.Equal(t, "init_action", st.States[2].TxnName, "Wrong state received")
	assert.Equal(t, "transaction_1", st.States[3].TxnName, "Wrong state received")
	assert.Equal(t, "transaction_2", st.States[4].TxnName, "Wrong state received")
	assert.Equal(t, "test_finish", st.States[5].TxnName, "Wrong state received")
	assert.Equal(t, "check_complete", st.States[6].TxnName, "Wrong state received")
	assert.Equal(t, "test_finish", st.States[7].TxnName, "Wrong state received")
	log.Printf("states: %+v", st.States)
	assert.Equal(t, false, p.State.Variables["success"])
	assert.Equal(t, false, p.State.Variables["failure"])
}

func TestRunPlanCorrectWithRemove(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	cc := make(chan bool)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "guid": "GUID", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				FailWithMessage(t, "Couldn't read body: "+err.Error())
			}
			var r interface{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				FailWithMessage(t, "couldn't marshal json: "+err.Error())
			}
			m := r.(map[string]interface{})
			s, ok := m["Michael"]
			if !ok {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if s != "Scott" {
				FailWithMessage(t, "JSON import failed")
				return
			}
			log.Printf("Callback succeeded!")
		} else if r.URL.Path == "/GUID" {
			cc <- true
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(finishedsuccessfulresultoneci))
			log.Printf("Callback succeeded!")
		} else {
			FailWithMessage(t, "Got incorrect URL for callback")
			return
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	//core.S.RunnerKillSwitch = true
	//time.Sleep(2 * time.Second)
	//core.S.RunnerKillSwitch = false

	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, resp := MakeRequest(o, "POST", "/capi/v1/launch/basic_test", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	p := handler.GetPlan()
	_, _ = MakeRequest(o, "POST", "/api/v1/url1", `{"amount": "too much", "I": "love tacos"}`, t, 200)
	_, _ = MakeRequest(o, "POST", "/api/v1/url2", "oompa loompa doopity doo", t, 200)

	// it's just too fast.
	//assert.Equal(t, "transaction_2", p.State.Transaction, "wth?")
	assert.Equal(t, "completed", p.State.States[2].Status, "Incorrect status %s", p.State.States[2].Status)
	assert.Equal(t, "expected", p.State.States[3].Status, "Incorrect status %s", p.State.States[3].Status)
	assert.Equal(t, "expected", p.State.States[4].Status, "Incorrect status %s", p.State.States[4].Status)
	assert.Equal(t, "init_action", p.State.States[2].TxnName, "Incorrect Transaction %s", p.State.States[2].TxnName)
	assert.Equal(t, "transaction_1", p.State.States[3].TxnName, "Incorrect Transaction %s", p.State.States[3].TxnName)
	assert.Equal(t, "transaction_2", p.State.States[4].TxnName, "Incorrect Transaction %s", p.State.States[4].TxnName)

	_ = <-cc

	_, _ = MakeRequest(o, "POST", "/capi/v1/remove", "", t, 200)

	// give the reset time to happen
	time.Sleep(200 * time.Millisecond)

	_, resp = MakeRequest(o, "POST", "/capi/v1/status", "", t, 200)
	body, _ := ioutil.ReadAll(resp.Body)
	st := api.HTTPReturnStruct{}
	err := json.Unmarshal(body, &st)
	if err != nil {
		FailWithMessage(t, err.Error())
		return
	}

	assert.Equal(t, "no change", st.Message)
}

func TestRunPlanIncorrectURLs(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				FailWithMessage(t, "Couldn't read body: "+err.Error())
			}
			var r interface{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				FailWithMessage(t, "couldn't marshal json: "+err.Error())
			}
			m := r.(map[string]interface{})
			s, ok := m["Michael"]
			if !ok {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if s != "Scott" {
				FailWithMessage(t, "JSON import failed")
				return
			}
			log.Printf("Callback succeeded!")
		} else {
			FailWithMessage(t, "Got incorrect URL for callback")
			return
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/basic_test", "", t, 200)

	p := handler.GetPlan()
	_, _ = MakeRequest(o, "POST", "/api/v1/url1a", `{"amount": "too much", "I": "love tacos"}`, t, 401)
	time.Sleep(10 * time.Second)

	assert.Equal(t, "transaction_1", p.State.Transaction, "Wasn't supposed to advance")
	assert.Equal(t, "unexpected", p.State.States[3].Status, "Incorrect status %s", p.State.States[3].Status)
	assert.Equal(t, "transaction_1", p.State.States[3].TxnName, "Incorrect Transaction %s", p.State.States[3].TxnName)
	assert.Equal(t, 4, len(p.State.States), "Incorrect length of State")
}

func TestRunPlanCorrectURLsIncorrectText(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				FailWithMessage(t, "Couldn't read body: "+err.Error())
			}
			var r interface{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				FailWithMessage(t, "couldn't marshal json: "+err.Error())
			}
			m := r.(map[string]interface{})
			s, ok := m["Michael"]
			if !ok {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if s != "Scott" {
				FailWithMessage(t, "JSON import failed")
				return
			}
			log.Printf("Callback succeeded!")
		} else {
			FailWithMessage(t, "Got incorrect URL for callback")
			return
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/basic_test", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	p := handler.GetPlan()

	_, _ = MakeRequest(o, "POST", "/api/v1/url1", `{"amount": "no such thing as too much", "I": "love tacos"}`, t, 401)
	assert.Equal(t, "transaction_1", p.State.Transaction, "Wasn't supposed to advance")
	assert.Equal(t, "unexpected", p.State.States[3].Status, "Incorrect status %s", p.State.States[3].Status)
	assert.Equal(t, "transaction_1", p.State.States[3].TxnName, "Incorrect Transaction %s", p.State.States[3].TxnName)
	assert.Equal(t, 4, len(p.State.States), "Incorrect length of State")
}

func TestFindConfigDir(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		want       string
	}{
		{
			name:       "Relative",
			configFile: "data/config.yml",
			want:       "data",
		},
		{
			name:       "Absolute",
			configFile: "/configs/config.yml",
			want:       "/configs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := o.ConfigFile
			o.ConfigFile = tt.configFile
			got := FindConfigDir(o)
			assert.Equal(t, tt.want, got)
			o.ConfigFile = old
		})
	}
}

func TestGoToConfigDir(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Valid",
			wantErr: false,
		},
		{
			name:    "Invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				configDir string
				err       error
			)
			assert := assert.New(t)
			configDir, err = ioutil.TempDir("", "config")
			if err != nil {
				t.Fatalf("Could not create temp dir: %s", err)
			}
			startDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Cannot determine working dir: %s", err)
			}
			if !tt.wantErr {
				defer os.RemoveAll(configDir)
				defer os.Chdir(startDir)
			} else {
				log.Printf("removing %s", configDir)
				err := os.RemoveAll(configDir)
				if err != nil {
					t.Errorf("couldn't remove %s", configDir)
				}
			}

			assert.NotEqual(configDir, startDir, "Starting WD")

			configFile := filepath.Join(configDir, "config.yml")

			old := o.ConfigFile
			o.ConfigFile = configFile
			log.Printf("configfile is %s", o.ConfigFile)
			err = GoToConfigDir(o)
			o.ConfigFile = old

			if tt.wantErr {
				assert.Error(err, "wanted error")
			} else {
				assert.NoError(err, "expected no error")
				finalDir, err := os.Getwd()
				if err != nil {
					t.Fatalf("Unable to determine final working dir: %s", err)
					assert.Equal(configDir, finalDir)
				}
			}
		})
	}
}

func TestRunPlanFailedIgnore(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testignorefailure", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()

	assert.Equal(t, "empty", p.State.Transaction, "State didn't advance")
}

func TestRunPlanFailedNonIgnore(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testnonignorefailure", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "ignorefailure_callback", p.State.Transaction, "State didn't advance")
}

func TestExternalVariables(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testextvariables", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	p := handler.GetPlan()
	assert.Equal(t, "somethingelse", p.State.Variables["variable2"], "variable2 not loaded")
}

func TestCounters(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testcounters", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()

	assert.Equal(t, 2.0, p.State.Variables["counter"])
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestDivide(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testdivide", "", t, 200)

	time.Sleep(5 * time.Second)
	p := handler.GetPlan()

	assert.Equal(t, 2.0, p.State.Variables["counter"])
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestMultiply(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testmultiply", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()

	time.Sleep(3 * time.Second)

	assert.Equal(t, 50.0, p.State.Variables["counter"])
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestSubtract(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testsubtract", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, 5.0, p.State.Variables["counter"])
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestGt(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testgt", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestLt(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testlt", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	time.Sleep(3 * time.Second)

	assert.Equal(t, "empty", p.State.Transaction)
}

func TestLe(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testle", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestGe(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testge", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	time.Sleep(3 * time.Second)

	assert.Equal(t, "empty", p.State.Transaction)
}

func TestEq(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testeq", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestNe(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testne", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestTransactionSubstitution(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/transactionvars", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction)
}

func TestSource(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/test_source", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction)
	assert.Equal(t, 1, p.State.Variables["destination"])
}

func TestURLVariable(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testurlvariable", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	_, _ = MakeRequest(o, "POST", "/api/v1/url1", `{"amount": "too much", "I": "love tacos"}`, t, 200)

	time.Sleep(3 * time.Second)

	_, ok := p.State.Variables["bert"]
	if !ok {
		FailWithMessage(t, fmt.Sprintf("Missing expected variable. %+v", p.State.Variables))
	}
	assert.Equal(t, `{"amount": "too much", "I": "love tacos"}`, p.State.Variables["bert"])

	_, ok = p.State.Variables["ernie"]
	if !ok {
		FailWithMessage(t, fmt.Sprintf("Missing expected variable. %+v", p.State.Variables))
	}
	assert.Equal(t, map[string]interface{}{"amount": "too much", "I": "love tacos"}, p.State.Variables["ernie"])
	// assert.Equal(t, "empty", State.Transaction.Name)
	// assert.Equal(t, 1, State.Plan.Variables["destination"])
}

func TestRunPlanJRAVA(t *testing.T) {
	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"map1": {"node1": "value", "node2": "value"}}`))
			log.Printf("Got callback %s", r.URL.Path)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				FailWithMessage(t, "Couldn't read body: "+err.Error())
			}
			var r interface{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				FailWithMessage(t, "couldn't marshal json: "+err.Error())
			}
			m := r.(map[string]interface{})
			s, ok := m["Michael"]
			if !ok {
				FailWithMessage(t, "JSON import failed")
				return
			}
			if s != "Scott" {
				FailWithMessage(t, "JSON import failed")
				return
			}
			log.Printf("Callback succeeded!")
		} else {
			FailWithMessage(t, "Got incorrect URL for callback")
			return
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/jrava", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	p := handler.GetPlan()
	time.Sleep(2 * time.Second)
	assert.Equal(t, "success", p.State.Transaction, "Didn't start in the right state")
}

func TestSatisfyGroupUrl1(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	//core.S.RunnerKillSwitch = true
	//time.Sleep(1 * time.Second)
	//core.S.RunnerKillSwitch = false

	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/satisfy_group", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	_, rresp := MakeRequest(o, "POST", "/api/v1/url1", `{"amount": "too much", "I": "love tacos"}`, t, 200)

	assert.Equal(t, "I love tacos", rresp.Body.String())
}

func TestSatisfyGroupUrl2(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	//core.S.RunnerKillSwitch = true
	//time.Sleep(1 * time.Second)
	//core.S.RunnerKillSwitch = false

	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/satisfy_group", "", t, 200)

	time.Sleep(200 * time.Millisecond)

	_, rresp := MakeRequest(o, "POST", "/api/v1/url2", `{"amount": "too much", "I": "love tacos"}`, t, 200)

	assert.Equal(t, "Don't quote me on this", rresp.Body.String())
}

func TestHeaders(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			log.Printf("header: %+v", r.Header)
			assert.Equal(t, "aheader1", r.Header.Get("X-Header1"))
			assert.Equal(t, "aheader2", r.Header.Get("X-Header2"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/header_test", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction, "State didn't advance")
}

func TestHeadersInvalid(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			log.Printf("header: %+v", r.Header)
			assert.Equal(t, "aheader1", r.Header.Get("X-Header1"))
			assert.Equal(t, "aheader2", r.Header.Get("X-Header2"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/header_testinvalid", "", t, 400)

}

func TestSplitCB(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			log.Printf("header: %+v", r.Header)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
		} else {
			t.Errorf("wrong path %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/cbsplit", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction, "State didn't advance")
}

func TestSplitCBWithWait(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			log.Printf("header: %+v", r.Header)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
		} else {
			t.Errorf("wrong path %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/cbsplitwait", "", t, 200)

	time.Sleep(10 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction, "State didn't advance")
}

func TestSplitCBWithWaitExpectingFailure(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			log.Printf("header: %+v", r.Header)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
		} else {
			t.Errorf("wrong path %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/cbsplitwaitfailure", "", t, 200)

	time.Sleep(10 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction, "State didn't advance")
}

func TestCallbackStringResponse(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			log.Printf("header: %+v", r.Header)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`Pong`))
			log.Printf("Got callback %s", r.URL.Path)
		} else {
			t.Errorf("wrong path %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/cbstring", "", t, 200)

	time.Sleep(10 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction, "State didn't advance")
}

func TestMatchString(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/testmatchstring", "", t, 200)

	time.Sleep(5 * time.Second)

	p := handler.GetPlan()
	assert.Equal(t, "empty", p.State.Transaction, "State didn't advance")
}

func TestSplitCBWithReset(t *testing.T) {

	InitHandler(t)
	defer StopHandler(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blah" {
			log.Printf("header: %+v", r.Header)
			time.Sleep(8 * time.Second)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message": "Something", "error": false}`))
			log.Printf("Got callback %s", r.URL.Path)
		} else {
			t.Errorf("wrong path %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	Config.Bases["testurl"] = ts.URL
	defer func(d string) { os.Chdir(d) }(workingDir(t))

	_, _ = MakeRequest(o, "POST", "/capi/v1/launch/cbsplitnofinish", "", t, 200)

	time.Sleep(2 * time.Second)

	_, _ = MakeRequest(o, "POST", "/capi/v1/remove", "", t, 200)

	time.Sleep(6 * time.Second)

	p := handler.GetPlan()
	// if plan is nil, this means the reset worked.  Maybe not properly - that is muuuch harder to figure out,
	// but worked.
	assert.Nil(t, p)
}
