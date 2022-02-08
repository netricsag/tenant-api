package util

import (
	"context"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	Clientset      *kubernetes.Clientset
	DISCOUNT_LABEL string
)

// GetPodsByTenant returns a map of pods for each tenant
func GetPodsByTenant(tenants []string) (map[string][]string, error) {
	tenantPods := make(map[string][]string)
	// get namespace with same name as tenant and get pods
	for _, tenant := range tenants {
		pods, err := Clientset.CoreV1().Pods(tenant).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		// for each pod add it to the list of pods for the namespace
		tenantPods[tenant] = make([]string, 0)
		for _, pod := range pods.Items {
			tenantPods[tenant] = append(tenantPods[tenant], pod.Name)
		}
	}

	return tenantPods, nil
}

// GetCPURequestsSumByTenant returns the sum of CPU requests for each tenant
func GetCPURequestsSumByTenant(tenants []string) (map[string]int64, error) {
	tenantCPURequests := make(map[string]int64)
	for _, tenant := range tenants {
		pods, err := Clientset.CoreV1().Pods(tenant).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			return nil, err
		}

		for _, pod := range pods.Items {

			// get DISCOUNT_REQUEST by DISCOUNT_LABEL
			discount := pod.Labels[DISCOUNT_LABEL]
			// convert to float64
			discountFloat, err := strconv.ParseFloat(discount, 64)
			if err != nil || discountFloat < 0 || discountFloat > 1 {
				return nil, err
			}

			CPU_DISCOUNT_PERCENT = discountFloat

			tenantCPURequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
		}
	}
	return tenantCPURequests, nil
}

// GetMemoryRequestsSumByTenant returns the sum of memory requests for each tenant
func GetMemoryRequestsSumByTenant(tenants []string) (map[string]int64, error) {
	tenantMemoryRequests := make(map[string]int64)
	for _, tenant := range tenants {
		pods, err := Clientset.CoreV1().Pods(tenant).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, pod := range pods.Items {

			// get DISCOUNT_REQUEST by DISCOUNT_LABEL
			discount := pod.Labels[DISCOUNT_LABEL]
			// convert to float64
			discountFloat, err := strconv.ParseFloat(discount, 64)
			if err != nil || discountFloat < 0 || discountFloat > 1 {
				return nil, err
			}

			MEMORY_DISCOUNT_PERCENT = discountFloat

			tenantMemoryRequests[tenant] += pod.Spec.Containers[0].Resources.Requests.Memory().Value()
		}
	}
	return tenantMemoryRequests, nil
}

// GetStorageRequestsSumByTenant returns the sum of storage requests for each tenant
func GetStorageRequestsSumByTenant(tenants []string) (map[string]map[string]int64, error) {
	tenantPVCs := make(map[string]map[string]int64)
	for _, tenant := range tenants {
		pvcList, err := Clientset.CoreV1().PersistentVolumeClaims(tenant).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			return nil, err
		}

		// create a map for each storage class with a count of pvc size if it exists
		tenantPVCs[tenant] = make(map[string]int64)
		for _, pvc := range pvcList.Items {
			// get DISCOUNT_REQUEST by DISCOUNT_LABEL
			discount := pvc.Labels[DISCOUNT_LABEL]
			// convert to float64
			discountFloat, err := strconv.ParseFloat(discount, 64)
			if err != nil || discountFloat < 0 || discountFloat > 1 {
				return nil, err
			}

			STORAGE_DISCOUNT_PERCENT = discountFloat
			tenantPVCs[tenant][*pvc.Spec.StorageClassName] += pvc.Spec.Resources.Requests.Storage().Value()
		}

		// if tenant is emtpy remove it from the map
		if len(tenantPVCs[tenant]) == 0 {
			delete(tenantPVCs, tenant)
		}
	}
	return tenantPVCs, nil
}

// GetIngressRequestsSumByTenant returns the sum of ingress requests for each tenant
func GetIngressRequestsSumByTenant(tenants []string) (map[string][]string, error) {
	tenantsIngress := make(map[string][]string)

	for _, tenant := range tenants {
		// get ingress for each namespace in the tenant and add it to the map of ingress for the tenant
		ingressList, err := Clientset.NetworkingV1().Ingresses(tenant).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			return nil, err
		}

		for _, ingress := range ingressList.Items {
			// get DISCOUNT_REQUEST by DISCOUNT_LABEL
			discount := ingress.Labels[DISCOUNT_LABEL]
			// convert to float64
			discountFloat, err := strconv.ParseFloat(discount, 64)
			if err != nil || discountFloat < 0 || discountFloat > 1 {
				return nil, err
			}

			INGRESS_DISCOUNT_PERCENT = discountFloat

			// apend ingress hostname to the list of ingress for the tenant
			tenantsIngress[tenant] = append(tenantsIngress[tenant], ingress.Name)
		}
	}

	return tenantsIngress, nil
}

// GetRessourceQuota returns the resource quota for the given tenant and label set in the config namespace
func GetRessourceQuota(tenant string, namespace_suffix string, label string) (float64, error) {
	// get the namespace with tenant-namespace_suffix and get the label value
	namespace, err := Clientset.CoreV1().Namespaces().Get(context.TODO(), tenant, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}

	// get the cpu quota from the label
	cpuQuota := namespace.Labels[label]
	if cpuQuota == "" {
		cpuQuota = "0"
	}

	// convert to float64
	cpuQuotaFloat, err := strconv.ParseFloat(cpuQuota, 64)
	if err != nil || cpuQuotaFloat < 0 {
		WarningLogger.Printf("CPU quota value %s is not valid for pod %s with label %s", cpuQuota, namespace.Name, label)
		cpuQuota = "0"
	}

	return cpuQuotaFloat, nil
}
