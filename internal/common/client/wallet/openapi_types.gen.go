// Package wallet provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package wallet

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	Code    float32 `json:"code"`
	Message string  `json:"message"`
}

// WalletInfo defines model for WalletInfo.
type WalletInfo struct {
	Address *string `json:"address,omitempty"`
	Chain   *string `json:"chain,omitempty"`
	Uuid    *string `json:"uuid,omitempty"`
}

// WalletParams defines model for WalletParams.
type WalletParams struct {
	Chain string `json:"chain"`
}

// CreateWalletJSONRequestBody defines body for CreateWallet for application/json ContentType.
type CreateWalletJSONRequestBody = WalletParams
