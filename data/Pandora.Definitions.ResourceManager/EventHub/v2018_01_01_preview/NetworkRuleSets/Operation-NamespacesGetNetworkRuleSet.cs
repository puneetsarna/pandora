using Pandora.Definitions.Operations;
using System.Collections.Generic;
using System.Net;
using Pandora.Definitions.Interfaces;

namespace Pandora.Definitions.ResourceManager.EventHub.v2018_01_01_preview.NetworkRuleSets
{
	internal class NamespacesGetNetworkRuleSet : GetOperation
	{
		public override ResourceID? ResourceId()
		{
			return new NamespaceId();
		}

		public override object? ResponseObject()
		{
			return new NetworkRuleSet();
		}

		public override string? UriSuffix()
		{
			return "/networkRuleSets/default";
		}
	}
}
