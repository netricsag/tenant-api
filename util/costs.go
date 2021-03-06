package util

import (
	"fmt"
	"strings"
)

var (
	CPU_COST                 float64
	MEMORY_COST              float64
	STORAGE_COST             map[string]map[string]float64
	INGRESS_COST             float64
	INGRESS_COST_PER_DOMAIN  bool
	EXCLUDE_INGRESS_VCLUSTER bool
	CPU_DISCOUNT_PERCENT     float64
	MEMORY_DISCOUNT_PERCENT  float64
	STORAGE_DISCOUNT_PERCENT float64
	INGRESS_DISCOUNT_PERCENT float64
)

// GetCPUCost returns the cost of the provided MiliCPU
func GetCPUCost(millicores float64) float64 {
	// return per core
	return (CPU_COST * float64(millicores) / 1000) * (1 - CPU_DISCOUNT_PERCENT)
}

// GetMemoryCost returns the cost of the provided Memory
func GetMemoryCost(memory float64) float64 {
	// return per GB
	return (MEMORY_COST * float64(memory) / (1024 * 1024 * 1024)) * (1 - MEMORY_DISCOUNT_PERCENT)
}

// GetStorageCost returns the cost of the provided Storage of the StorageClass
func GetStorageCost(storageClass string, size float64) (float64, error) {
	// return per GB
	if STORAGE_COST[storageClass] == nil {
		return 0, fmt.Errorf("storage class %s not found", storageClass)
	}
	// with STORAGE_DISCOUNT_PERCENT
	return (STORAGE_COST[storageClass]["cost"] * float64(size) / (1024 * 1024 * 1024)) * (1 - STORAGE_DISCOUNT_PERCENT), nil
}

// GetIngressCost returns the cost of the provided Ingress
func GetIngressCostByDomain(hostnameStrings []string) float64 {

	var tenantIngressCostsPerDomainSum float64

	// define a set of domains
	domains := make(map[string]bool)

	for _, host := range hostnameStrings {
		// split string with .
		hostnameParts := strings.Split(host, ".")
		// get the 2 last parts of the hostname
		var domain string
		if len(hostnameParts) > 1 {
			domain = hostnameParts[len(hostnameParts)-2] + "." + hostnameParts[len(hostnameParts)-1]
		} else {
			domain = ""
		}
		// add the domain to the tenantIngressCosts map
		if domain != "" {
			domains[domain] = true
		} else {
			ErrorLogger.Printf("domain is not valid for hostname %s", host)
		}
	}

	// calculate the cost * count of domains
	for range domains {
		tenantIngressCostsPerDomainSum += INGRESS_COST * (1 - INGRESS_DISCOUNT_PERCENT)
	}

	return tenantIngressCostsPerDomainSum
}

func GetIngressCost(ingressCount int) float64 {
	return INGRESS_COST * float64(ingressCount) * (1 - INGRESS_DISCOUNT_PERCENT)
}
