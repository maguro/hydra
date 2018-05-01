/*
 * ORY Hydra - Cloud Native OAuth 2.0 and OpenID Connect Server
 *
 * Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here. Keep in mind that this document reflects the latest branch, always. Support for versioned documentation is coming in the future.
 *
 * OpenAPI spec version: Latest
 * Contact: hi@ory.am
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package swagger

type ConsentRequest struct {

	// Challenge is the identifier (\"authorization challenge\") of the consent authorization request. It is used to identify the session.
	Challenge string `json:"challenge,omitempty"`

	Client OAuth2Client `json:"client,omitempty"`

	OidcContext OpenIdConnectContext `json:"oidc_context,omitempty"`

	// RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters.
	RequestUrl string `json:"request_url,omitempty"`

	// RequestedScope contains all scopes requested by the OAuth 2.0 client.
	RequestedScope []string `json:"requested_scope,omitempty"`

	// Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the consent request using the usual API call.
	Skip bool `json:"skip,omitempty"`

	// Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client.
	Subject string `json:"subject,omitempty"`
}
