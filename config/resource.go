// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	NativeCurrencyAddr = "0x0000000000000000000000000000000000000000"
)

type resourceConfigFile struct {
	Resources []Resource `json:"resources"`
}

type resourceDomainInfo struct {
	DomainId int `json:"domainId"`
	Decimal  int `json:"decimal"`
}

type Resource struct {
	ID      string               `json:"id"`
	Symbol  string               `json:"symbol"`
	Domains []resourceDomainInfo `json:"domains"`
}

func newResource(resourceID string, symbol string, domains []resourceDomainInfo) *Resource {
	return &Resource{
		ID:      resourceID,
		Symbol:  symbol,
		Domains: domains,
	}
}

// loadResources registers and load all pre-defined resources
func loadResources(resourceConfigPath string) map[string]Resource {
	resourceData, err := ioutil.ReadFile(filepath.Clean(resourceConfigPath))
	if err != nil {
		panic(ErrLoadResourceConfig.Wrap(err))
	}

	return parseResources(resourceData)
}

// ResourceIDBuilder builds the resourceID according to fee handler contract
func ResourceIDBuilder(tokenAddress string, domainId int) string {
	return fmt.Sprintf("%s%d", strings.ToLower(tokenAddress), domainId)
}

func parseResources(resourceData []byte) map[string]Resource {
	var content resourceConfigFile
	err := json.Unmarshal(resourceData, &content)
	if err != nil {
		panic(ErrLoadResourceConfig.Wrap(err))
	}

	resources := make(map[string]Resource, 0)
	for _, resource := range content.Resources {
		resources[strings.ToLower(resource.ID)] =
			*newResource(resource.ID, resource.Symbol, resource.Domains)
	}

	return resources
}
