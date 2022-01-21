/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/natron-io/tenant-api/controllers"
	"github.com/natron-io/tenant-api/routes"
	"github.com/natron-io/tenant-api/util"

	"github.com/joho/godotenv"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func init() {
	util.InitLoggers()

	if err := godotenv.Load(); err != nil {
		util.WarningLogger.Println("Error loading .env file")
	}

	if util.CLIENT_ID = os.Getenv("CLIENT_ID"); util.CLIENT_ID == "" {
		util.WarningLogger.Println("CLIENT_ID is not set")
	}

	if util.CLIENT_SECRET = os.Getenv("CLIENT_SECRET"); util.CLIENT_SECRET == "" {
		util.WarningLogger.Println("CLIENT_SECRET is not set")
	}

	if controllers.SECRET_KEY = os.Getenv("SECRET_KEY"); controllers.SECRET_KEY == "" {
		util.WarningLogger.Println("SECRET_KEY is not set")
		// setting random key
		controllers.SECRET_KEY = util.RandomStringBytes(32)
		util.InfoLogger.Printf("SECRET_KEY is not set, using random key: %s", controllers.SECRET_KEY)
	}

	if controllers.LABELSELECTOR = os.Getenv("LABELSELECTOR"); controllers.LABELSELECTOR == "" {
		util.WarningLogger.Println("LABELSELECTOR is not set")
		controllers.LABELSELECTOR = "natron.io/tenant"
		util.InfoLogger.Printf("LABELSELECTOR set using default: %s", controllers.LABELSELECTOR)
	}

	if controllers.CALLBACK_URL = os.Getenv("CALLBACK_URL"); controllers.CALLBACK_URL == "" {
		util.WarningLogger.Println("CALLBACK_URL is not set")
		controllers.CALLBACK_URL = "http://localhost:3000"
		util.InfoLogger.Printf("CALLBACK_URL set using default: %s", controllers.CALLBACK_URL)
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	util.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET",
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		// set header to html
		c.Set("Content-Type", "text/html")
		// return html for login to github
		return c.SendString("<a href='/login/github'>Login with Github</a>")
	})

	routes.Setup(app, util.Clientset)

	app.Listen(":8000")

	util.InfoLogger.Println("Tenant API is running on port 8000")

	// for {
	// 	// get pods in all the namespaces by omitting namespace
	// 	// Or specify namespace to get pods in particular namespace
	// 	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// 	// Examples for error handling:
	// 	// - Use helper functions e.g. errors.IsNotFound()
	// 	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	// 	_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
	// 	if errors.IsNotFound(err) {
	// 		fmt.Printf("Pod example-xxxxx not found in default namespace\n")
	// 	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 		fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	// 	} else if err != nil {
	// 		panic(err.Error())
	// 	} else {
	// 		fmt.Printf("Found example-xxxxx pod in default namespace\n")
	// 	}

	// 	time.Sleep(10 * time.Second)
	// }
}
