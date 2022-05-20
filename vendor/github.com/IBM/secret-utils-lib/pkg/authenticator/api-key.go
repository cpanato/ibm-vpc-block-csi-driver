/**
 * Copyright 2022 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authenticator

import (
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secret-utils-lib/pkg/token"
	"github.com/IBM/secret-utils-lib/pkg/utils"
	"go.uber.org/zap"
)

// APIKeyAuthenticator ...
type APIKeyAuthenticator struct {
	authenticator     *core.IamAuthenticator
	logger            *zap.Logger
	isSecretEncrypted bool
}

// NewIamAuthenticator ...
func NewIamAuthenticator(apikey string, logger *zap.Logger) *APIKeyAuthenticator {
	logger.Info("Initializing iam authenticator")
	defer logger.Info("Initialized iam authenticator")
	aa := new(APIKeyAuthenticator)
	aa.authenticator = new(core.IamAuthenticator)
	aa.authenticator.ApiKey = apikey
	aa.logger = logger
	return aa
}

// GetToken ...
func (aa *APIKeyAuthenticator) GetToken(freshTokenRequired bool) (string, uint64, error) {
	aa.logger.Info("Fetching IAM token using api key authenticator")
	var iamtoken string
	var err error
	var tokenlifetime uint64

	if !freshTokenRequired {
		aa.logger.Info("Request received to fetch existing token")
		iamtoken, err = aa.authenticator.GetToken()
		if err != nil {
			aa.logger.Error("Error fetching existing token", zap.Error(err))
			return "", tokenlifetime, utils.Error{Description: "Error fetching iam token using api key", BackendError: err.Error()}
		}

		// Fetching token life time
		tokenlifetime, err = token.CheckTokenLifeTime(iamtoken)
		if err == nil {
			aa.logger.Info("Fetched iam token and token lifetime successfully")
			return iamtoken, tokenlifetime, nil
		}
		aa.logger.Error("Error fetching token lifetime of existing token", zap.Error(err))
	}

	aa.logger.Info("Fetching fresh token")
	tokenResponse, err := aa.authenticator.RequestToken()
	if err != nil {
		aa.logger.Error("Error fetching fresh token", zap.Error(err))
		return "", tokenlifetime, utils.Error{Description: "Error fetching iam token using api key", BackendError: err.Error()}
	}

	if tokenResponse == nil {
		aa.logger.Error("Token response received is empty")
		return "", tokenlifetime, utils.Error{Description: utils.ErrEmptyTokenResponse}
	}

	iamtoken = tokenResponse.AccessToken
	tokenlifetime, err = token.CheckTokenLifeTime(iamtoken)
	if err != nil {
		aa.logger.Error("Error fetching token lifetime for new token", zap.Error(err))
		return "", tokenlifetime, utils.Error{Description: "Error fetching token lifetime", BackendError: err.Error()}
	}

	aa.logger.Info("Successfully fetched IAM token and token lifetime")
	return iamtoken, tokenlifetime, nil
}

// GetSecret ...
func (aa *APIKeyAuthenticator) GetSecret() string {
	return aa.authenticator.ApiKey
}

// SetSecret ...
func (aa *APIKeyAuthenticator) SetSecret(secret string) {
	aa.authenticator.ApiKey = secret
}

// SetURL ...
func (aa *APIKeyAuthenticator) SetURL(url string) {
	aa.authenticator.URL = url
}

// IsSecretEncrypted ...
func (aa *APIKeyAuthenticator) IsSecretEncrypted() bool {
	return aa.isSecretEncrypted
}

// SetEncryption ...
func (aa *APIKeyAuthenticator) SetEncryption(encrypted bool) {
	aa.isSecretEncrypted = encrypted
}