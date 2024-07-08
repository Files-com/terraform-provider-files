package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	providerConfig = `
provider "files" {
  api_key = "test"
}
`
)

var clientsLock = sync.RWMutex{}
var clients = map[string]*http.Client{}

func VcrTest(t *testing.T, c resource.TestCase) {
	defer closeRecorder(t)
	resource.UnitTest(t, c)
}

func ProviderFactories(testName string) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"files": providerserver.NewProtocol6WithError(newTestProvider(testName)),
	}
}

func newTestProvider(testName string) *testProvider {
	return &testProvider{
		filesProvider: filesProvider{version: "test"},
		testName:      testName,
	}
}

type testProvider struct {
	filesProvider
	testName string
}

func (t *testProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	t.filesProvider.Configure(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	client := getCachedClient(t.testName, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	config := files_sdk.Config{}.Init().SetCustomClient(client)

	resp.DataSourceData = config
	resp.ResourceData = config
}

func getCachedClient(testName string, resp *provider.ConfigureResponse) *http.Client {
	clientsLock.RLock()
	client, ok := clients[testName]
	clientsLock.RUnlock()
	if ok {
		return client
	}

	var r *recorder.Recorder
	var err error
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		r, err = recorder.NewAsMode(filepath.Join("fixtures", testName), recorder.ModeReplaying, nil)
	} else {
		r, err = recorder.New(filepath.Join("fixtures", testName))
	}
	if err != nil {
		resp.Diagnostics.AddError("recorder.New", err.Error())
		return nil
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Key")
		return nil
	})

	r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		if cassette.DefaultMatcher(r, i) {
			if r.Body != nil {
				io.ReadAll(r.Body)
				r.Body.Close()
			}

			return true
		}
		return false
	})

	client = &http.Client{Transport: r}
	clientsLock.Lock()
	clients[testName] = client
	clientsLock.Unlock()
	return client
}

func closeRecorder(t *testing.T) {
	clientsLock.RLock()
	client, ok := clients[t.Name()]
	clientsLock.RUnlock()
	if ok {
		if !t.Failed() {
			err := client.Transport.(*recorder.Recorder).Stop()
			if err != nil {
				t.Error(err)
			}
		}

		clientsLock.Lock()
		delete(clients, t.Name())
		clientsLock.Unlock()
	}
}
