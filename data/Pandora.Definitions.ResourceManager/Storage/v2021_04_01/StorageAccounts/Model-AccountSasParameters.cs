using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;

namespace Pandora.Definitions.ResourceManager.Storage.v2021_04_01.StorageAccounts;


internal class AccountSasParametersModel
{
    [JsonPropertyName("keyToSign")]
    public string? KeyToSign { get; set; }

    [DateFormat(DateFormatAttribute.DateFormat.RFC3339)]
    [JsonPropertyName("signedExpiry")]
    [Required]
    public DateTime SignedExpiry { get; set; }

    [JsonPropertyName("signedIp")]
    public string? SignedIp { get; set; }

    [JsonPropertyName("signedPermission")]
    [Required]
    public PermissionsConstant SignedPermission { get; set; }

    [JsonPropertyName("signedProtocol")]
    public HttpProtocolConstant? SignedProtocol { get; set; }

    [JsonPropertyName("signedResourceTypes")]
    [Required]
    public SignedResourceTypesConstant SignedResourceTypes { get; set; }

    [JsonPropertyName("signedServices")]
    [Required]
    public ServicesConstant SignedServices { get; set; }

    [DateFormat(DateFormatAttribute.DateFormat.RFC3339)]
    [JsonPropertyName("signedStart")]
    public DateTime? SignedStart { get; set; }
}
