using Pandora.Definitions.Attributes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System;
using System.Collections.Generic;
using System.Net;

namespace Pandora.Definitions.ResourceManager.ManagedIdentity.v2018_11_30.ManagedIdentity
{
    internal class UserAssignedIdentitiesGetOperation : Operations.GetOperation
    {
        public override ResourceID? ResourceId() => new UserAssignedIdentitiesId();

        public override Type? ResponseObject() => typeof(IdentityModel);


    }
}
