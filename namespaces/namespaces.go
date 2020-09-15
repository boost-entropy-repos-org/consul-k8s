// Package namespaces handles interaction with Consul namespaces needed across
// commands.
package namespaces

import (
	"fmt"

	capi "github.com/hashicorp/consul/api"
)

// EnsureExists ensures a Consul namespace with name ns exists. If it doesn't,
// it will create it and set crossNSACLPolicy as a policy default.
func EnsureExists(client *capi.Client, ns string, crossNSAClPolicy string) error {
	// Check if the Consul namespace exists.
	namespaceInfo, _, err := client.Namespaces().Read(ns, nil)
	if err != nil {
		return err
	}
	if namespaceInfo != nil {
		return nil
	}

	// If not, create it.
	var aclConfig capi.NamespaceACLConfig
	if crossNSAClPolicy != "" {
		// Create the ACLs config for the cross-Consul-namespace
		// default policy that needs to be attached
		aclConfig = capi.NamespaceACLConfig{
			PolicyDefaults: []capi.ACLLink{
				{Name: crossNSAClPolicy},
			},
		}
	}

	consulNamespace := capi.Namespace{
		Name:        ns,
		Description: "Auto-generated by consul-k8s",
		ACLs:        &aclConfig,
		Meta:        map[string]string{"external-source": "kubernetes"},
	}

	_, _, err = client.Namespaces().Create(&consulNamespace, nil)
	return err
}

// ConsulNamespace returns the consul namespace that a service should be
// registered in based on the namespace options. It returns an
// empty string if namespaces aren't enabled.
func ConsulNamespace(kubeNS string, enableConsulNamespaces bool, consulDestNS string, enableMirroring bool, mirroringPrefix string) string {
	if !enableConsulNamespaces {
		return ""
	}

	// Mirroring takes precedence.
	if enableMirroring {
		return fmt.Sprintf("%s%s", mirroringPrefix, kubeNS)
	}

	return consulDestNS
}