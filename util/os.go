package util

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	err  error
	CORS string
)

// LoadEnv loads OS environment variables
func LoadEnv() error {
	if CLIENT_ID = os.Getenv("CLIENT_ID"); CLIENT_ID == "" {
		err = errors.New("CLIENT_ID is not set")
		ErrorLogger.Println(err)
		Status = "Error: CLIENT_ID is not set"
		return err
	}

	if CLIENT_SECRET = os.Getenv("CLIENT_SECRET"); CLIENT_SECRET == "" {
		err = errors.New("CLIENT_SECRET is not set")
		ErrorLogger.Println(err)
		Status = "Error: CLIENT_SECRET is not set"
		return err
	}

	if CALLBACK_URL = os.Getenv("CALLBACK_URL"); CALLBACK_URL == "" {
		WarningLogger.Println("CALLBACK_URL is not set")
		CALLBACK_URL = "http://localhost:3000"
		InfoLogger.Printf("CALLBACK_URL set using default: %s", CALLBACK_URL)
	} else {
		InfoLogger.Printf("CALLBACK_URL set using env: %s", CALLBACK_URL)
	}

	if CORS = os.Getenv("CORS"); CORS == "" {
		WarningLogger.Println("CORS is not set")
		CORS = "*"
		InfoLogger.Printf("CORS set using default: %s", CORS)
	} else {
		InfoLogger.Printf("CORS set using env: %s", CORS)
	}

	if MAX_REQUESTS, err = strconv.Atoi(os.Getenv("MAX_REQUESTS")); err != nil {
		WarningLogger.Println("MAX_REQUESTS is not set")
		MAX_REQUESTS = 100
		InfoLogger.Printf("MAX_REQUESTS set using default: %d", MAX_REQUESTS)
	} else {
		InfoLogger.Printf("MAX_REQUESTS set using env: %d", MAX_REQUESTS)
	}

	if SECRET_KEY = os.Getenv("SECRET_KEY"); SECRET_KEY == "" {
		WarningLogger.Println("SECRET_KEY is not set")
		// setting random key
		SECRET_KEY = RandomStringBytes(32)
		InfoLogger.Printf("SECRET_KEY is not set, using random key: %s", SECRET_KEY)
	}

	if DISCOUNT_LABEL = os.Getenv("DISCOUNT_LABEL"); DISCOUNT_LABEL == "" {
		WarningLogger.Println("DISCOUNT_LABEL is not set")
		DISCOUNT_LABEL = "natron.io/discount"
		InfoLogger.Printf("DISCOUNT_LABEL set using default: %s", DISCOUNT_LABEL)
	} else {
		InfoLogger.Printf("DISCOUNT_LABEL set using env: %s", DISCOUNT_LABEL)
	}

	if CPU_COST, err = strconv.ParseFloat(os.Getenv("CPU_COST"), 64); CPU_COST == 0 || err != nil {
		WarningLogger.Println("CPU_COST is not set or invalid float value")
		CPU_COST = 1.00
		InfoLogger.Printf("CPU_COST set using default: %f", CPU_COST)
	} else {
		InfoLogger.Printf("CPU_COST set using env: %f", CPU_COST)
	}

	if MEMORY_COST, err = strconv.ParseFloat(os.Getenv("MEMORY_COST"), 64); MEMORY_COST == 0 || err != nil {
		WarningLogger.Println("MEMORY_COST is not set or invalid float value")
		MEMORY_COST = 1.00
		InfoLogger.Printf("MEMORY_COST set using default: %f", MEMORY_COST)
	} else {
		InfoLogger.Printf("MEMORY_COST set using env: %f", MEMORY_COST)
	}

	if INGRESS_COST, err = strconv.ParseFloat(os.Getenv("INGRESS_COST"), 64); INGRESS_COST == 0 || err != nil {
		WarningLogger.Println("INGRESS_COST is not set or invalid float value")
		INGRESS_COST = 1.00
		InfoLogger.Printf("INGRESS_COST set using default: %f", INGRESS_COST)
	} else {
		InfoLogger.Printf("INGRESS_COST set to: %f", INGRESS_COST)
	}

	if INGRESS_COST_PER_DOMAIN, err = strconv.ParseBool(os.Getenv("INGRESS_COST_PER_DOMAIN")); !INGRESS_COST_PER_DOMAIN || err != nil {
		WarningLogger.Println("INGRESS_COST_PER_DOMAIN is not set or invalid bool value")
		INGRESS_COST_PER_DOMAIN = false
		InfoLogger.Printf("INGRESS_COST_PER_DOMAIN set using default: %t", INGRESS_COST_PER_DOMAIN)
	} else {
		InfoLogger.Printf("INGRESS_COST_PER_DOMAIN set using env: %t", INGRESS_COST_PER_DOMAIN)
	}

	if EXCLUDE_INGRESS_VCLUSTER, err = strconv.ParseBool(os.Getenv("EXCLUDE_INGRESS_VCLUSTER")); !EXCLUDE_INGRESS_VCLUSTER || err != nil {
		WarningLogger.Println("EXCLUDE_INGRESS_VCLUSTER is not set or invalid bool value")
		EXCLUDE_INGRESS_VCLUSTER = false
		InfoLogger.Printf("EXCLUDE_INGRESS_VCLUSTER set using default: %t", EXCLUDE_INGRESS_VCLUSTER)
	} else {
		InfoLogger.Printf("EXCLUDE_INGRESS_VCLUSTER set using env: %t", EXCLUDE_INGRESS_VCLUSTER)
	}

	if SLACK_TOKEN = os.Getenv("SLACK_TOKEN"); SLACK_TOKEN == "" {
		WarningLogger.Println("SLACK_TOKEN is not set")
		SLACK_TOKEN = ""
	} else {
		InfoLogger.Printf("SLACK_TOKEN set using env: %s", SLACK_TOKEN)
	}

	if BroadCastChannelID = os.Getenv("SLACK_BROADCAST_CHANNEL_ID"); BroadCastChannelID == "" && SLACK_TOKEN != "" {
		ErrorLogger.Println("SLACK_BROADCAST_CHANNEL_ID is not set")
		Status = "Error: SLACK_BROADCAST_CHANNEL_ID is not set"
		os.Exit(1)
	} else {
		InfoLogger.Printf("SLACK_BROADCAST_CHANNEL_ID set using env: %s", BroadCastChannelID)
	}

	if SlackURL = os.Getenv("SLACK_URL"); SlackURL == "" && SLACK_TOKEN != "" {
		ErrorLogger.Println("SLACK_URL is not set")
		Status = "Error: SLACK_URL is not set"
		os.Exit(1)
	} else {
		InfoLogger.Printf("SLACK_URL set using env: %s", SlackURL)
	}

	// ======================== //
	// 		StorageClasses		//
	// ======================== //
	// get every env variable starting with STORAGE_COST_ and parse it to STORAGE_COST with the storage class name after STORAGE_COST_ as key
	tempStorageCost := make(map[string]map[string]float64)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "STORAGE_COST_") {
			// split env variable to key and value
			keyValue := strings.Split(env, "=")
			// split key to storage class name and cost
			key := strings.Split(keyValue[0], "_")
			// parse value to float
			value, err := strconv.ParseFloat(keyValue[1], 64)
			if err != nil {
				err = errors.New("STORAGE_COST_" + key[2] + " is not set or invalid float value")
				ErrorLogger.Println(err)
				Status = "Error: " + err.Error()
				return err
			}
			// add to tempStorageCost
			tempStorageCost[key[2]] = map[string]float64{"cost": value}
			InfoLogger.Printf("storage class %s set to cost value: %f", key[2], value)
		}
	}
	STORAGE_COST = tempStorageCost

	storageClassesInCluster, err := GetStorageClassesInCluster()
	if err != nil {
		err = errors.New("cannot get storage classes in cluster")
		ErrorLogger.Println(err)
		Status = "Error: " + err.Error()
		os.Exit(1)
	}

	// check if every storage class in cluster is in STORAGE_COST
	for _, storageClass := range storageClassesInCluster {
		if _, ok := STORAGE_COST[storageClass]; !ok {
			err = errors.New("Storage class " + storageClass + " is not set")
			ErrorLogger.Println(err)
			Status = "Error: " + err.Error()
			os.Exit(1)
		}
	}

	if STORAGE_COST == nil {
		WarningLogger.Println("STORAGE_COST is not set")
		STORAGE_COST = map[string]map[string]float64{
			"default": {"cost": 1.00},
		}
		InfoLogger.Printf("cost for storage class default set using default: %f", STORAGE_COST["default"]["cost"])
	}

	return nil
}
