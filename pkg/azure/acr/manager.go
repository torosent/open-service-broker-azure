package acr

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/containerregistry"
	"github.com/Azure/go-autorest/autorest/azure"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
)

// Manager is an interface to be implemented by any component capable of
// managing an acr instance
type Manager interface {
	DeleteServer(
		registryName string,
		resourceGroupName string,
	) error
}

type manager struct {
	azureEnvironment azure.Environment
	subscriptionID   string
	tenantID         string
	clientID         string
	clientSecret     string
}

// NewManager returns a new implementation of the Manager interface
func NewManager() (Manager, error) {
	azureConfig, err := az.GetConfig()
	if err != nil {
		return nil, err
	}
	azureEnvironment, err := azure.EnvironmentFromName(azureConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf(
			`error parsing Azure environment name "%s"`,
			azureConfig.Environment,
		)
	}
	return &manager{
		azureEnvironment: azureEnvironment,
		subscriptionID:   azureConfig.SubscriptionID,
		tenantID:         azureConfig.TenantID,
		clientID:         azureConfig.ClientID,
		clientSecret:     azureConfig.ClientSecret,
	}, nil
}

func (m *manager) DeleteServer(
	registryName string,
	resourceGroupName string,
) error {
	authorizer, err := az.GetBearerTokenAuthorizer(
		m.azureEnvironment,
		m.tenantID,
		m.clientID,
		m.clientSecret,
	)
	if err != nil {
		return fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	servicesClient := containerregistry.NewRegistriesClientWithBaseURI(
		m.azureEnvironment.ResourceManagerEndpoint,
		m.subscriptionID,
	)
	servicesClient.Authorizer = authorizer
	_, err = servicesClient.Delete(
		resourceGroupName,
		registryName,
	)
	if err != nil {
		return fmt.Errorf("error deleting Azure acr: %s", err)
	}

	return nil
}
