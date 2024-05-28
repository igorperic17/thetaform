// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-theta/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name scaffolding

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/theta. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "hashicorp.com/edu/theta",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

// func main() {

// 	plugin.Serve(&plugin.ServeOpts{
// 		ProviderFunc: func() *schema.Provider {
// 			return &schema.Provider{
// 				Schema: map[string]*schema.Schema{
// 					"base_url": {
// 						Type:        schema.TypeString,
// 						Required:    true,
// 						DefaultFunc: schema.EnvDefaultFunc("THETA_BASE_URL", nil),
// 					},
// 					"api_key": {
// 						Type:        schema.TypeString,
// 						Required:    true,
// 						DefaultFunc: schema.EnvDefaultFunc("THETA_API_KEY", nil),
// 					},
// 				},
// 				ResourcesMap: map[string]*schema.Resource{
// 					"theta_endpoint": provider.NewEndpoint(),
// 				},
// 				ConfigureFunc: configureProvider,
// 			}
// 		},
// 	})
// }

// func configureProvider(d *schema.ResourceData) (interface{}, error) {
// 	baseURL := d.Get("base_url").(string)
// 	apiKey := d.Get("api_key").(string)
// 	return provider.NewClient(baseURL, apiKey), nil
// }