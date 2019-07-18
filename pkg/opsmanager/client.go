// Package opsmanager hosts a HTTP client which abstracts communication with Ops Manager instances.
//
// To create a new client, you have to call the following code:
//
//		resolver := httpclient.NewURLResolverWithPrefix("http://OPS-MANAGER-INSTANCE", opsmanager.PublicApiPrefix)
// 		client := opsmanager.NewClient(WithResolver(resolver))
//
// The client can then be used to issue requests such as:
//
//		user := User{Username: ..., Password: ..., ...}
// 		globalOwner := client.CreateFirstUser(user, WhitelistAllowAll)
//
// See the Client interface below for a list of all the support operations.
// If however, you need one that is not currently supported, the _opsManagerApi_ struct extends
// _httpclient.BasicHTTPOperation_, allowing you to issue raw HTTP requests to the specified Ops Manager instance.
//
// 		url := resolver.Of("/path/to/a/resource/%s", id)
//		resp:= client.GetJSON(url)
//		useful.PanicOnUnrecoverableError(resp.Err)
//		defer useful.LogError(resp.Response.Body.Close)
//		var data SomeType
//		decoder := json.NewDecoder(resp.Response.Body)
//		err := decoder.Decode(&result)
//		useful.PanicOnUnrecoverableError(err)
//		// do something with _data_
//
// To issue authenticated requests, initialize a client using:
//		withResolver := opsmanager.WithResolver(httpclient.NewURLResolverWithPrefix("http://OPS-MANAGER-INSTANCE", PublicAPIPrefix))
//		withDigestAuth := httpclient.WithDigestAuthentication(publicKey, privateKey)
//		withHTTPClient := opsmanager.WithHTTPClient(httpclient.NewClient(withDigestAuth))
//		client := opsmanager.NewClient(withResolver, withHTTPClient)
//
// The following can be used for authentication:
//		- Ops Manager user credentials: (username, password) - only works with some APIs and should not be used
//		- Programmatic API keys: (publicKey, privateKey) - preferred credentials pair
//		- Ops Manager user and a Personal API Key: (username, personalAPIKey) - deprecated
// You can read more about this topic here: https://docs.opsmanager.mongodb.com/master/tutorial/configure-public-api-access/#configure-public-api-access
//
package opsmanager

import (
	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/useful"
	"github.com/pkg/errors"
	"io"
)

type opsManagerClient struct {
	httpclient.BasicClient

	resolver httpclient.URLResolver
}

// Client defines the API actions implemented in this client
type Client interface {
	httpclient.BasicClient

	// https://docs.opsmanager.mongodb.com/master/reference/api/user-create-first/
	CreateFirstUser(user User, whitelistIP string) (CreateFirstUserResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/groups/get-all-groups-for-current-user/
	GetAllProjects() (ProjectsResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/groups/create-one-group/
	CreateOneProject(name string, orgID string) (CreateOneProjectResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/#get-the-automation-configuration
	GetAutomationConfig(projectID string) (AutomationConfig, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/agents-get-by-type/
	GetAgentsByType(projectID string, agentType string) (GetAgentsByTypeResponse, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/automation-config/#update-the-automation-configuration
	UpdateDeployments(projectID string, body io.Reader) (map[string]interface{}, error)
	// https://docs.opsmanager.mongodb.com/master/reference/api/agentapikeys/create-one-agent-api-key/
	CreateAgentAPIKEY(projectID string, desc string) (CreateAgentAPIKEYResponse, error)
}

// NewClient builds a new API client for connecting to Ops Manager
func NewClient(configs ...func(*opsManagerClient)) Client {
	// initialize a bare client
	client := &opsManagerClient{}

	// apply all configurations
	for _, configure := range configs {
		configure(client)
	}

	// validations
	if client.resolver == nil {
		useful.PanicOnUnrecoverableError(errors.New("the client requires a URLResolver with the appropriate Ops Manager URL configured"))
	}
	if client.BasicClient == nil {
		useful.PanicOnUnrecoverableError(errors.New("the client requires an underlying basic HTTP client to be configured"))
	}

	return client
}

// WithResolver configures an Ops Manager client which relies on the specified resolver
func WithResolver(resolver httpclient.URLResolver) func(*opsManagerClient) {
	return func(client *opsManagerClient) {
		client.resolver = resolver
	}
}

// WithHTTPClient configures an Ops Manager which delegates basic HTTP operations to the specified client
func WithHTTPClient(basicClient httpclient.BasicClient) func(*opsManagerClient) {
	return func(client *opsManagerClient) {
		client.BasicClient = basicClient
	}
}