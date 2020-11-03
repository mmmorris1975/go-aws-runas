package credentials

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/mmmorris1975/aws-runas/credentials/helpers"
	"github.com/mmmorris1975/aws-runas/shared"
	"os"
	"time"
)

type baseStsProvider struct {
	credentials.Expiry
	Client        stsiface.STSAPI
	Cache         CredentialCacher
	Duration      time.Duration
	ExpiryWindow  time.Duration
	Logger        shared.Logger
	SerialNumber  string
	TokenCode     string
	TokenProvider func() (string, error)
}

func newBaseStsProvider(cfg client.ConfigProvider) *baseStsProvider {
	return &baseStsProvider{
		Client:        sts.New(cfg),
		Logger:        new(shared.DefaultLogger),
		TokenProvider: helpers.NewMfaTokenProvider(os.Stdin).ReadInput,
	}
}

// CheckCache will load credentials from cache.  If a cache is not configured, this method will
// return an empty and expired set of credentials.
func (p *baseStsProvider) CheckCache() *Credentials {
	var creds *Credentials

	if p.Cache != nil {
		if creds = p.Cache.Load(); creds.Value().HasKeys() {
			p.Logger.Debugf("loaded sts credentials from cache")
			p.SetExpiration(creds.Expiration, 0)
		} else {
			p.SetExpiration(time.Unix(0, 0), 0)
		}
	}

	return creds
}

// ConvertDuration normalizes and returns an int64 duration value which is compatible with the AWS
// SDK credential duration field in the API input objects.  The 1st duration argument to this method
// will be checked against the other provided duration values.  If less than 1, the default value will
// be used, if less than the minimum, the minimum value will be used, and if greater than the maximum,
// the maximum value will be used.
func (p *baseStsProvider) ConvertDuration(d, min, max, def time.Duration) *int64 {
	if d < 1 {
		p.Logger.Debugf("provided duration less than 1, setting to default value")
		d = def
	} else if d < min {
		p.Logger.Debugf("provided duration too short, setting to minimum value")
		d = min
	} else if d > max {
		p.Logger.Debugf("provided duration too long, setting to maximum value")
		d = max
	}

	return aws.Int64(int64(d.Seconds()))
}

func (p *baseStsProvider) handleMfa() (*string, error) {
	if len(p.SerialNumber) > 0 {
		if len(p.TokenCode) > 0 {
			return aws.String(p.TokenCode), nil
		}

		// prompt for mfa
		if p.TokenProvider != nil {
			t, err := p.TokenProvider()
			if err != nil {
				return nil, err
			}
			return aws.String(t), nil
		}

		return nil, ErrMfaRequired
	}

	// mfa not required
	return nil, nil
}