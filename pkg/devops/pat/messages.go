package pat

type PatToken struct {
	AuthorizationId string   `json:"authorizationId"`
	DisplayName     string   `json:"displayName"`
	Scope           string   `json:"scope"`
	TargetAccounts  []string `json:"targetAccounts"`
	Token           string   `json:"token"`
	ValidFrom       string   `json:"validFrom"`
	ValidTo         string   `json:"validTo"`
}

type PagedPatTokens struct {
	ContinuationToken string     `json:"ContinuationToken"`
	PatTokens         []PatToken `json:"patTokens"`
}

type PatTokenCreateRequest struct {
	AllOrgs     bool   `json:"allOrgs"`
	DisplayName string `json:"displayName"`
	Scope       string `json:"scope"`
	ValidTo     string `json:"validTo"`
}

type PatTokenUpdateRequest struct {
	AllOrgs         bool   `json:"allOrgs"`
	AuthorizationId string `json:"authorizationId"`
	DisplayName     string `json:"displayName"`
	Scope           string `json:"scope"`
	ValidTo         string `json:"validTo"`
}

type PatTokenResult struct {
	PatToken      PatToken `json:"patToken"`
	PatTokenError string   `json:"patTokenError"`
}

const (
	AccessDenied                = "accessDenied"
	AuthorizationNotFound       = "authorizationNotFound"
	DisplayNameRequired         = "displayNameRequired"
	DuplicateHash               = "duplicateHash"
	FailedToIssueAccessToken    = "failedToIssueAccessToken"
	FailedToReadTenantPolicy    = "failedToReadTenantPolicy"
	FailedToUpdateAccessToken   = "failedToUpdateAccessToken"
	FullScopePatPolicyViolation = "fullScopePatPolicyViolation"
	GlobalPatPolicyViolation    = "globalPatPolicyViolation"
	HostAuthorizationNotFound   = "hostAuthorizationNotFound"
	InvalidAuthorizationId      = "invalidAuthorizationId"
	InvalidClient               = "invalidClient"
	InvalidClientId             = "invalidClientId"
	InvalidClientType           = "invalidClientType"
	InvalidDisplayName          = "invalidDisplayName"
	InvalidScope                = "invalidScope"
	InvalidSource               = "invalidSource"
	InvalidSourceIP             = "invalidSourceIP"
	InvalidTargetAccounts       = "invalidTargetAccounts"
	InvalidToken                = "invalidToken"
	InvalidUserId               = "invalidUserId"
	InvalidUserType             = "invalidUserType"
	InvalidValidTo              = "invalidValidTo"
	None                        = "none"
	PatLifespanPolicyViolation  = "patLifespanPolicyViolation"
	SourceNotSupported          = "sourceNotSupported"
	SshPolicyDisabled           = "sshPolicyDisabled"
	TokenNotFound               = "tokenNotFound"
	UserIdRequired              = "userIdRequired"
)
