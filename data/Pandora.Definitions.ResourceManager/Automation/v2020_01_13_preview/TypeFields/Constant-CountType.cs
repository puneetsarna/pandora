using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.Automation.v2020_01_13_preview.TypeFields;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum CountTypeConstant
{
    [Description("nodeconfiguration")]
    Nodeconfiguration,

    [Description("status")]
    Status,
}
