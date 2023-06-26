using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.Consumption.v2023_05_01.ReservationRecommendationDetails;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum TermConstant
{
    [Description("P1Y")]
    POneY,

    [Description("P3Y")]
    PThreeY,
}
