using Pandora.Definitions.Attributes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System.Collections.Generic;
using System.Net;

namespace Pandora.Definitions.ResourceManager.Relay.v2017_04_01.Namespaces
{
    internal class UpdateOperation : Operations.PatchOperation
    {
        public override object? RequestObject()
        {
            return new RelayUpdateParametersModel();
        }

        public override ResourceID? ResourceId()
        {
            return new NamespaceId();
        }

        public override object? ResponseObject()
        {
            return new RelayNamespaceModel();
        }


    }
}
