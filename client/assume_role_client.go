/*
 * Copyright (c) 2021 Michael Morris. All Rights Reserved.
 *
 * Licensed under the MIT license (the "License"). You may not use this file except in compliance
 * with the License. A copy of the License is located at
 *
 * https://github.com/mmmorris1975/aws-runas/blob/master/LICENSE
 *
 * or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License
 * for the specific language governing permissions and limitations under the License.
 */

package client

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/mmmorris1975/aws-runas/credentials"
	"os/user"
	"time"
)

type assumeRoleClient struct {
	*baseIamClient
	provider *credentials.AssumeRoleProvider
}

// AssumeRoleClientConfig is the configuration attributes for the STS Assume Role operation for either IAM identities,
// or role chaining using SAML or OIDC Identity Tokens.
type AssumeRoleClientConfig struct {
	SessionTokenClientConfig
	RoleArn         string
	RoleSessionName string
	ExternalId      string
}

// NewAssumeRoleClient is an AwsClient which knows how to do Assume Role operations.
func NewAssumeRoleClient(cfg aws.Config, clientCfg *AssumeRoleClientConfig) *assumeRoleClient {
	c := &assumeRoleClient{newBaseIamClient(cfg, clientCfg.Logger), nil}

	p := credentials.NewAssumeRoleProvider(cfg, clientCfg.RoleArn)
	p.Cache = clientCfg.Cache
	p.Duration = clientCfg.Duration
	p.SerialNumber = clientCfg.SerialNumber
	p.TokenCode = clientCfg.TokenCode
	p.TokenProvider = clientCfg.TokenProvider
	p.ExternalId = clientCfg.ExternalId
	p.RoleSessionName = clientCfg.RoleSessionName
	p.Logger = clientCfg.Logger

	if len(p.RoleSessionName) < 2 { // AWS SDK minimum length
		if id, err := c.ident.Identity(); err == nil {
			p.RoleSessionName = id.Username
		} else if usr, err := user.Current(); err == nil {
			p.RoleSessionName = usr.Username
		} else {
			// escape route value ... matches AWS SDK value defaulting logic
			p.RoleSessionName = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
		}
	}

	c.provider = p
	c.creds = aws.NewCredentialsCache(p, func(o *aws.CredentialsCacheOptions) {
		o.ExpiryWindow = p.ExpiryWindow
	})
	return c
}

// ClearCache cleans the cache for this client's AWS credential cache.
func (c *assumeRoleClient) ClearCache() error {
	if c.creds != nil {
		c.creds.Invalidate()
	}

	if c.provider.Cache != nil {
		c.provider.Logger.Debugf("clearing cached assume role credentials")
		return c.provider.Cache.Clear()
	}
	return nil
}
