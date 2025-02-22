// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	v1 "github.com/hashicorp/pandora/tools/data-api-sdk/v1"
	"github.com/hashicorp/pandora/tools/data-api-sdk/v1/models"
	"github.com/hashicorp/pandora/tools/generator-go-sdk/internal/generator"
	"github.com/mitchellh/cli"
)

var _ cli.Command = GenerateCommand{}

type GenerateCommand struct {
	sourceDataType models.SourceDataType
}

type GeneratorInput struct {
	apiServerEndpoint string
	outputDirectory   string
	services          []string
	settings          generator.Settings
}

func NewGenerateCommand(sourceDataType models.SourceDataType) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return GenerateCommand{
			sourceDataType: sourceDataType,
		}, nil
	}
}

func (g GenerateCommand) Help() string {
	// TODO: expand with commands
	return "Generates a Go SDK based on the API Definitions from the Data API"
}

func (g GenerateCommand) Run(args []string) int {
	ctx := context.Background()

	input := GeneratorInput{
		settings: generator.Settings{},
	}

	input.settings.UseOldBaseLayerFor(
		// @tombuildsstuff: New Services should now use the `hashicorp/go-azure-sdk` base layer by default
		// instead of the base layer from `Azure/go-autorest` - as such this list is for compatibility purposes
		// with services already used in `terraform-provider-azurerm`. These services will be gradually removed
		// from this list to ensure they're migrated across to using `hashicorp/go-azure-sdk`s base layer.

		"ContainerInstance",
		"CosmosDB",
		"FrontDoor",
		"Insights",
		"RecoveryServicesBackup", // error: generating Service "RecoveryServicesBackup" / Version "2023-04-01" / Resource "Operation": generating methods: templating methods (using hashicorp/go-azure-sdk): templating: building methods: building response struct template: existing model "ValidateOperationResponse" conflicts with the operation response model for "Validate"
		"Security",
		"SecurityInsights",
		"ServiceFabric",
		"Subscription",

		// @tombuildsstuff: The Key Vault API has an issue where it requires that the EXACT casing returned in the Response
		// is sent in the Request to update or remove a Key Vault Access Policy - and using other casings mean the update
		// or removal fails - which is tracked in https://github.com/hashicorp/pandora/issues/3229.
		//
		// After testing it appears that `2023-07-01` doesn't suffer from this problem - as such we're going to leave
		// `2023-02-01` on the older base layer and use the newer API Version as a divide to give us a clear migration path.
		"KeyVault@2023-02-01",
	)

	var serviceNames string

	f := flag.NewFlagSet("generator-go-sdk", flag.ExitOnError)
	f.StringVar(&input.apiServerEndpoint, "data-api", "http://localhost:5000", "-data-api=http://localhost:5000")
	f.StringVar(&input.outputDirectory, "output-dir", "", "-output-dir=../generated-sdk-dev")
	f.StringVar(&serviceNames, "services", "", "A list of comma separated Service named from the Data API to import")
	if err := f.Parse(args); err != nil {
		log.Fatalf("parsing arguments: %+v", err)
	}

	if serviceNames != "" {
		input.services = strings.Split(serviceNames, ",")
	}

	if input.outputDirectory == "" {
		homeDir, _ := os.UserHomeDir()
		input.outputDirectory = filepath.Join(homeDir, "/Desktop/generated-sdk-dev")
	}

	if err := g.run(ctx, input); err != nil {
		log.Fatalf("running generator: %+v", err)
	}

	return 0
}

func (g GenerateCommand) Synopsis() string {
	return "Generates a Go SDK based on the API Definitions from the Data API"
}

func (g GenerateCommand) run(ctx context.Context, input GeneratorInput) error {
	// output into a directory named after the source data type (e.g. `{dir}/resource-manager`)
	input.outputDirectory = path.Join(input.outputDirectory, string(g.sourceDataType))

	client := v1.NewClient(input.apiServerEndpoint, g.sourceDataType)
	data, err := client.LoadAllData(ctx, input.services)
	if err != nil {
		return fmt.Errorf("retrieving API Definitions: %+v", err)
	}

	errCh := make(chan error, 1)
	waitDone := make(chan struct{}, 1)
	var wg sync.WaitGroup
	addErr := func(err error) {
		// only put one err to channel
		select {
		case errCh <- err:
		default:
		}
	}

	generatorService := generator.NewServiceGenerator(input.settings)
	for serviceName, service := range data.Services {
		log.Printf("[DEBUG] Service %q..", serviceName)
		if !service.Generate {
			log.Printf("[DEBUG] .. is opted out of generation, skipping..")
			continue
		}

		wg.Add(1)
		go func(serviceName string, service models.Service, input GeneratorInput) {
			defer wg.Done()
			log.Printf("[DEBUG] Service %q", serviceName)
			for versionNumber, versionDetails := range service.APIVersions {
				log.Printf("[DEBUG]   Version %q", versionNumber)
				for resourceName, resourceDetails := range versionDetails.Resources {
					log.Printf("[DEBUG]      Resource %q..", resourceName)
					generatorData := generator.ServiceGeneratorInput{
						ServiceName:     serviceName,
						ServiceDetails:  service,
						VersionName:     versionNumber,
						VersionDetails:  versionDetails,
						ResourceName:    resourceName,
						ResourceDetails: resourceDetails,
						OutputDirectory: input.outputDirectory,
						Source:          versionDetails.Source,
					}
					log.Printf("[DEBUG] Generating Service %q / Version %q / Resource %q..", serviceName, versionNumber, resourceName)
					if err := generatorService.Generate(generatorData); err != nil {
						addErr(fmt.Errorf("generating Service %q / Version %q / Resource %q: %+v", serviceName, versionNumber, resourceName, err))
						return
					}
					log.Printf("[DEBUG] Generated Service %q / Version %q / Resource %q..", serviceName, versionNumber, resourceName)
				}

				// then output the Meta Client
				generatorData := generator.VersionInput{
					OutputDirectory: input.outputDirectory,
					ServiceName:     serviceName,
					VersionName:     versionNumber,
					Resources:       versionDetails.Resources,
					Source:          versionDetails.Source,
				}
				generatorData.UseNewBaseLayer = false
				if input.settings.ShouldUseNewBaseLayer(serviceName, versionNumber) {
					generatorData.UseNewBaseLayer = true
				}
				log.Printf("[DEBUG] Generating Service %q / Version %q..", serviceName, versionNumber)
				if err := generatorService.GenerateForVersion(generatorData); err != nil {
					addErr(fmt.Errorf("generating Service %q / Version %q: %+v", serviceName, versionNumber, err))
					return
				}
				log.Printf("[DEBUG] Generated Service %q / Version %q..", serviceName, versionNumber)
			}
		}(serviceName, service, input)
	}

	go func() {
		wg.Wait()
		waitDone <- struct{}{}
	}()

	select {
	case <-waitDone:
		break
	case err := <-errCh:
		return err
	}

	return nil
}
