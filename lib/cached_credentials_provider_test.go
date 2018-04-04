package lib

import (
	"github.com/mbndr/logo"
	"os"
	"strings"
	"testing"
	"time"
)

// common provider options
var opts = &CachedCredentialsProviderOptions{LogLevel: logo.DEBUG, cacheFilePrefix: "test"}

func TestNewCachedCredentialsProvider(t *testing.T) {
	t.Run("ProfileNil", func(t *testing.T) {
		defer func() {
			if x := recover(); x == nil {
				t.Errorf("Did not receive expected panic calling NewCachedCredentialsProvider with nil profile")
			}
		}()
		NewCachedCredentialsProvider(nil, opts)
	})

	t.Run("OptionsNil", func(t *testing.T) {
		p := NewCachedCredentialsProvider(new(AWSProfile), nil)
		if !strings.HasSuffix(p.cacher.CacheFile(), "_") {
			t.Errorf("Unexpected value returned calling NewCachedCredentialsProvider with nil options")
		}
	})
}

func TestCachedCredentialsProvider_CacheFile(t *testing.T) {
	t.Run("SourceProfileSet", func(t *testing.T) {
		prof := &AWSProfile{Name: "mock", SourceProfile: "source"}
		p := NewCachedCredentialsProvider(prof, opts)
		if !strings.HasSuffix(p.CacheFile(), ".aws_assume_role_source") {
			t.Errorf("Unexpected value for cache file name with profile name set: %s", p.CacheFile())
		}
	})

	t.Run("SourceProfileUnset", func(t *testing.T) {
		prof := &AWSProfile{Name: "mock", SourceProfile: ""}
		p := NewCachedCredentialsProvider(prof, opts)
		if !strings.HasSuffix(p.CacheFile(), ".aws_assume_role_mock") {
			t.Errorf("Unexpected value for cache file name with profile name set: %s", p.CacheFile())
		}
	})
}

func TestCachedCredentialsProvider_IsExpired(t *testing.T) {
	t.Run("CredsNil", func(t *testing.T) {
		p := NewCachedCredentialsProvider(new(AWSProfile), opts)
		if !p.IsExpired() {
			t.Errorf("Expected IsExpired() to be true for nil creds")
		}
	})

	t.Run("True", func(t *testing.T) {
		p := NewCachedCredentialsProvider(new(AWSProfile), opts)
		p.creds = &CachableCredentials{Expiration: 500}
		p.cacher = &credentialsCacher{file: os.DevNull}
		if !p.IsExpired() {
			t.Errorf("Expected IsExpired() to be true for expired creds")
		}
	})

	t.Run("False", func(t *testing.T) {
		p := NewCachedCredentialsProvider(new(AWSProfile), opts)
		p.creds = &CachableCredentials{Expiration: time.Now().Unix() + 500}
		p.cacher = &credentialsCacher{file: os.DevNull}
		if p.IsExpired() {
			t.Errorf("Expected IsExpired() to be false for non-expired creds")
		}
	})
}

func TestCachedCredentialsProvider_ExpirationTime(t *testing.T) {
	t.Run("CredsNil", func(t *testing.T) {
		p := NewCachedCredentialsProvider(new(AWSProfile), opts)
		if p.ExpirationTime() != time.Unix(0, 0) {
			t.Errorf("Expected nil credentials to have epoch expiration time, got :%v", p.ExpirationTime())
		}
	})

	t.Run("CredsValid", func(t *testing.T) {
		p := NewCachedCredentialsProvider(new(AWSProfile), opts)
		p.creds = &CachableCredentials{Expiration: time.Now().Unix()}
		if p.ExpirationTime() == time.Unix(0, 0) {
			t.Errorf("Expected valid credentials to not have epoch expiration time, got :%v", p.ExpirationTime())
		}
	})
}
