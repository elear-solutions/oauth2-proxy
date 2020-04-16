package providers

import (
	"fmt"
	"net/http"
	"net/url"
        "errors"
        "github.com/oauth2-proxy/oauth2-proxy/pkg/apis/sessions"
        "github.com/oauth2-proxy/oauth2-proxy/pkg/requests"
)

type COCOProvider struct {
	*ProviderData
}

// NewCOCOProvider initiates a new COCOProvider
func NewCOCOProvider(p *ProviderData) *COCOProvider {
	p.ProviderName = "COCO"
	if p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{Scheme: "https",
			Host: "api.staging.getcoco.buzz",
			Path: "/oauth/authorize",
			// ?granted_scopes=true
		}
	}
	if p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{Scheme: "https",
			Host: "api.staging.getcoco.buzz",
			Path: "/oauth/token",
		}
	}
	if p.ProfileURL.String() == "" {
		p.ProfileURL = &url.URL{Scheme: "https",
			Host: "api.staging.getcoco.buzz",
			Path: "/user-manager/users/me",
		}
	}
	if p.ValidateURL.String() == "" {
		p.ValidateURL = p.ProfileURL
	}
	if p.Scope == "" {
		p.Scope = "profile"
	}
	return &COCOProvider{ProviderData: p}
}

// segrigate coco header from acccess token

func getCOCOHeader(accessToken string) http.Header {
	header := make(http.Header)
	header.Set("Accept", "application/json")
	header.Set("x-li-format", "json")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
        fmt.Println("returnd header:")
	return header
}


// GetEmailAddress returns the Account userId 

func (p *COCOProvider) GetEmailAddress(s *sessions.SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ProfileURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header = getCOCOHeader(s.AccessToken)

	type result struct {
		Email string `json:"username"`
	}
	var r result
	err = requests.RequestJSON(req, &r)
	if err != nil {
		return "", err
	}
	if r.Email == "" {
		return "", errors.New("no email")
	}

	return r.Email , nil
}

// ValidateSessionState validates the AccessToken
func (p *COCOProvider) ValidateSessionState(s *sessions.SessionState) bool {
	return validateToken(p, s.AccessToken, getCOCOHeader(s.AccessToken))
}
