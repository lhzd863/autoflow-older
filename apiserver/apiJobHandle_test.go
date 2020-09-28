package apiserver

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/lhzd863/autoflow/module"
	"github.com/lhzd863/autoflow/util"
)

func Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/api").
		Consumes(restful.MIME_JSON, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_JSON)

	s := new(module.MetaApiServerBean)
	s.HomeDir = "/home/k8s/autoflow/autoflow-server/data/apiserver"
	r := NewResponseResourceJob(s)

	ws.Route(ws.POST("/status/pending").To(r.FlowJobStatusGetPendingHandle))

	container.Add(ws)
}

func RunRestfulCurlyRouterServer() {
	// setup service
	wsContainer := restful.NewContainer()
	Register(wsContainer)

	log.Print("start listening on localhost:12306")
	server := &http.Server{Addr: ":12306", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

func waitForServerUp(serverURL string) error {
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(5 * time.Second) {
		_, err := http.Get(serverURL + "/api/status/pending")
		if err == nil {
			return nil
		}
	}
	return errors.New("waiting for server timed out")
}

func TestFlowJobStatusGetPendingHandle(t *testing.T) {
	serverURL := "http://localhost:12306"
	go func() {
		RunRestfulCurlyRouterServer()
	}()

	if err := waitForServerUp(serverURL); err != nil {
		t.Errorf("%v", err)
	}

	LogF = "/home/k8s/autoflow/autoflow-server/data/apiserver/t.log"
	statusPendingHashRing = util.NewConsistent()
	for i := 0; i < 6; i++ {
		si := fmt.Sprintf("%d", i)
		statusPendingHashRing.Add(util.NewNode(i, si, 1))
	}
	statusGoHashRing = util.NewConsistent()
	for i := 0; i < 6; i++ {
		si := fmt.Sprintf("%d", i)
		statusGoHashRing.Add(util.NewNode(i, si, 1))
	}
	statusPendingOffset = 0
	statusGoOffset = 0

	// Send a POST request.
	var jsonStr = []byte(`{"id":"2020070113522949sdsdsa","flowid":"20200701135229494jre","ishash":"0","status":"Penidng"}`)
	req, err := http.NewRequest("POST", serverURL+"/api/status/pending", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected response: %v, expected: %v", resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	fmt.Println(string(body))
}

func setup() {
	fmt.Println("Before all tests")
}

func teardown() {
	fmt.Println("After all tests")
}

// func TestMain(m *testing.M) {
// 	setup()
// 	code := m.Run()
// 	teardown()
// 	os.Exit(code)
// }
