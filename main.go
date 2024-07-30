// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"
	"os"
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

	// var client = provider.NewClient("igorperic+theta2@live.com", "theta1231")

	// organisations, err := client.GetOrganizations()
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// }
	// project := organisations[0]
	// println(project.Name)

	// newProject := provider.Project{
	// 	Name: "my_cool_new_project",
	// }

	// projectResponse, err := client.CreateProject(&newProject)
	// if err != nil {
	// 	println(err.Error())
	// }

	// println(projectResponse.Name)

	// project := provider.Project{ID: "prj_2tcq81a5wyk3637caigcn9ytrz2s"}

	// templateReq := provider.DeploymentTemplateRequest{
	// 	Name:           "asdf2123",
	// 	ProjectID:      project.ID,
	// 	Description:    "",
	// 	ContainerImage: "vllm/vllm-openai",
	// 	Tags:           []string{"LLM", "API"}, // TODO: convert this into enums of allowed tags
	// 	ContainerPort:  "8000",
	// 	ContainerArgs:  "",
	// 	EnvVars:        nil,
	// }

	// template, err := client.CreateDeploymentTemplate(templateReq)
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// } else {
	// 	println("Created template")
	// 	println(template.ID)
	// }

	// templates, _ := client.GetDeploymentTemplates("prj_2tcq81a5wyk3637caigcn9ytrz2s", 0, 8)
	// println(templates[0].Name)

	// template, err := client.DeleteDeploymentTemplate("img_6xtzvg40d7cmbnpu1u13u58m6fn11", project.ID)
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// } else {
	// 	println("Deleted template")
	// 	println(template)
	// }

	// newTemplate, err := client.GetDeployment(template.ID)
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// } else {
	// 	println("Created template")
	// 	println(newTemplate)
	// }

	// Set log output to standard error for better visibility
	log.SetOutput(os.Stderr)
	// Set log level to debug for detailed logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("DEBUG: ")

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
