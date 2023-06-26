using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.Kusto.v2023_05_02.DataConnections;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum DatabaseRoutingConstant
{
    [Description("Multi")]
    Multi,

    [Description("Single")]
    Single,
}
