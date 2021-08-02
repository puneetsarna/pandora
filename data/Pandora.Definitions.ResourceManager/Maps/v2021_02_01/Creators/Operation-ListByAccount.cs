using Pandora.Definitions.Attributes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System.Collections.Generic;
using System.Net;

namespace Pandora.Definitions.ResourceManager.Maps.v2021_02_01.Creators
{
    internal class ListByAccountOperation : Operations.ListOperation
    {
        public override string? FieldContainingPaginationDetails()
        {
            return "nextLink";
        }

        public override ResourceID? ResourceId()
        {
            return new AccountId();
        }

        public override object NestedItemType()
        {
            return new CreatorModel();
        }

        public override string? UriSuffix()
        {
            return "/creators";
        }


    }
}
