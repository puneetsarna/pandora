using Pandora.Definitions.Attributes;
using Pandora.Definitions.CustomTypes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System;
using System.Collections.Generic;
using System.Net;

namespace Pandora.Definitions.ResourceManager.Storage.v2021_04_01.BlobService;

internal class GetServicePropertiesOperation : Operations.GetOperation
{
    public override ResourceID? ResourceId() => new BlobServiceId();

    public override Type? ResponseObject() => typeof(BlobServicePropertiesModel);


}
