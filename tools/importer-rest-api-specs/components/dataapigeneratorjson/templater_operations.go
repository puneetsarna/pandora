// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataapigeneratorjson

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	importerModels "github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
	"github.com/hashicorp/pandora/tools/sdk/dataapimodels"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func codeForOperation(operationName string, input importerModels.OperationDetails, knownConstants map[string]resourcemanager.ConstantDetails, knownModels map[string]importerModels.ModelDetails) ([]byte, error) {
	contentType := input.ContentType
	if strings.Contains(strings.ToLower(input.ContentType), "application/json") {
		contentType = fmt.Sprintf("%s; charset=utf-8", input.ContentType)
	}

	output := dataapimodels.Operation{
		Name:                             operationName,
		ContentType:                      contentType,
		ExpectedStatusCodes:              input.ExpectedStatusCodes,
		FieldContainingPaginationDetails: input.FieldContainingPaginationDetails,
		LongRunning:                      input.LongRunning,
		HTTPMethod:                       strings.ToUpper(input.Method),
		ResourceIdName:                   input.ResourceIdName,
		UriSuffix:                        input.UriSuffix,
	}

	if input.RequestObject != nil {
		requestObject, err := mapObjectDefinition(input.RequestObject, knownConstants, knownModels)
		if err != nil {
			return nil, fmt.Errorf("mapping the request object definition: %+v", err)
		}
		output.RequestObject = requestObject
	}

	if input.ResponseObject != nil {
		responseObject, err := mapObjectDefinition(input.ResponseObject, knownConstants, knownModels)
		if err != nil {
			return nil, fmt.Errorf("mapping the response object definition: %+v", err)
		}
		output.ResponseObject = responseObject
	}

	if len(input.Options) > 0 {
		options := make([]dataapimodels.Option, 0)
		sortedOptionsKeys := make([]string, 0)
		for k := range input.Options {
			sortedOptionsKeys = append(sortedOptionsKeys, k)
		}
		sort.Strings(sortedOptionsKeys)

		for _, optionName := range sortedOptionsKeys {
			optionDetails := input.Options[optionName]

			if optionDetails.ObjectDefinition == nil {
				return nil, fmt.Errorf("missing object definition")
			}

			optionObjectDefinition, err := mapOptionObjectDefinition(optionDetails.ObjectDefinition, knownConstants, knownModels)
			if err != nil {
				return nil, fmt.Errorf("mapping the object definition: %+v", err)
			}

			option := dataapimodels.Option{
				HeaderName:       optionDetails.HeaderName,
				QueryString:      optionDetails.QueryStringName,
				Field:            optionName,
				ObjectDefinition: optionObjectDefinition,
			}

			if !optionDetails.Required {
				option.Optional = true
			} else {
				option.Required = true
			}

			options = append(options, option)
		}
		output.Options = pointer.To(options)
	}

	data, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return nil, err
	}

	return data, nil
}

var internalObjectDefinitionsToOptionObjectDefinitionTypes = map[importerModels.ObjectDefinitionType]dataapimodels.OptionObjectDefinitionType{
	importerModels.ObjectDefinitionBoolean:   dataapimodels.BooleanOptionObjectDefinitionType,
	importerModels.ObjectDefinitionCsv:       dataapimodels.CsvOptionObjectDefinitionType,
	importerModels.ObjectDefinitionInteger:   dataapimodels.IntegerOptionObjectDefinitionType,
	importerModels.ObjectDefinitionFloat:     dataapimodels.FloatOptionObjectDefinitionType,
	importerModels.ObjectDefinitionList:      dataapimodels.ListOptionObjectDefinitionType,
	importerModels.ObjectDefinitionReference: dataapimodels.ReferenceOptionObjectDefinitionType,
	importerModels.ObjectDefinitionString:    dataapimodels.StringOptionObjectDefinitionType,
}

func mapOptionObjectDefinition(definition *importerModels.ObjectDefinition, constants map[string]resourcemanager.ConstantDetails, models map[string]importerModels.ModelDetails) (*dataapimodels.OptionObjectDefinition, error) {
	typeVal, ok := internalObjectDefinitionsToOptionObjectDefinitionTypes[definition.Type]
	if !ok {
		return nil, fmt.Errorf("internal-error: no OptionObjectDefinition mapping is defined for the OptionObjectDefinition Type %q", string(definition.Type))
	}

	output := dataapimodels.OptionObjectDefinition{
		Type:          typeVal,
		ReferenceName: nil,
		NestedItem:    nil,
	}

	if definition.ReferenceName != nil {
		output.ReferenceName = definition.ReferenceName
	}

	if definition.NestedItem != nil {
		nestedItem, err := mapOptionObjectDefinition(definition.NestedItem, constants, models)
		if err != nil {
			return nil, fmt.Errorf("mapping nested option object definition: %+v", err)
		}
		output.NestedItem = nestedItem
	}

	// finally let's do some sanity-checking to ensure the data being output looks legit
	if err := validateOptionObjectDefinition(output, constants, models); err != nil {
		return nil, fmt.Errorf("validating mapped OptionObjectDefinition: %+v", err)
	}

	return &output, nil
}

func validateOptionObjectDefinition(input dataapimodels.OptionObjectDefinition, constants map[string]resourcemanager.ConstantDetails, models map[string]importerModels.ModelDetails) error {
	requiresNestedItem := input.Type == dataapimodels.CsvOptionObjectDefinitionType || input.Type == dataapimodels.ListOptionObjectDefinitionType
	requiresReference := input.Type == dataapimodels.ReferenceOptionObjectDefinitionType
	if requiresNestedItem && input.NestedItem == nil {
		return fmt.Errorf("a Nested Object Definition must be specified for a %q type but didn't get one", string(input.Type))
	}
	if !requiresNestedItem && input.NestedItem != nil {
		return fmt.Errorf("a Nested Object Definition must not be specified for a %q type but got %q", string(input.Type), string(input.NestedItem.Type))
	}
	if requiresReference {
		if input.ReferenceName == nil {
			return fmt.Errorf("a Reference must be specified for a %q type but didn't get one", string(input.Type))
		}

		_, isConstant := constants[*input.ReferenceName]
		_, isModel := models[*input.ReferenceName]
		if !isConstant && !isModel {
			return fmt.Errorf("reference %q was not found as a constant or a model", *input.ReferenceName)
		}
		if isConstant && isModel {
			return fmt.Errorf("internal-error: %q was found as BOTH a Constant and a Model", *input.ReferenceName)
		}
	}
	if !requiresReference && input.ReferenceName != nil {
		return fmt.Errorf("a Reference must not be specified for a %q type but got %q", string(input.Type), *input.ReferenceName)
	}

	return nil
}
