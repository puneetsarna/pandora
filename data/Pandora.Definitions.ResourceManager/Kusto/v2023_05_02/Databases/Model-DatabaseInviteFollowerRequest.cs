using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.Kusto.v2023_05_02.Databases;


internal class DatabaseInviteFollowerRequestModel
{
    [JsonPropertyName("inviteeEmail")]
    [Required]
    public string InviteeEmail { get; set; }

    [JsonPropertyName("tableLevelSharingProperties")]
    public TableLevelSharingPropertiesModel? TableLevelSharingProperties { get; set; }
}
