using Pandora.Definitions.Attributes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System.Collections.Generic;
using System.Net;

namespace Pandora.Definitions.ResourceManager.AppConfiguration.v2020_06_01.PrivateLinkResources
{
    internal class ListByConfigurationStoreOperation : Operations.ListOperation
    {
        public override string? FieldContainingPaginationDetails()
        {
            return "nextLink";
        }

        public override ResourceID? ResourceId()
        {
            return new ConfigurationStoreId();
        }

        public override object NestedItemType()
        {
            return new PrivateLinkResourceModel();
        }

        public override string? UriSuffix()
        {
            return "/privateLinkResources";
        }


    }
}
