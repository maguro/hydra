/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
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
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent

import (
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

type MemoryManager struct {
	consentRequests        map[string]ConsentRequest
	handledConsentRequests map[string]HandledConsentRequest
	authRequests           map[string]AuthenticationRequest
	handledAuthRequests    map[string]HandledAuthenticationRequest
	authSessions           map[string]AuthenticationSession
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		consentRequests:        map[string]ConsentRequest{},
		handledConsentRequests: map[string]HandledConsentRequest{},
		authRequests:           map[string]AuthenticationRequest{},
		handledAuthRequests:    map[string]HandledAuthenticationRequest{},
		authSessions:           map[string]AuthenticationSession{},
	}
}

func (m *MemoryManager) CreateConsentRequest(c *ConsentRequest) error {
	m.consentRequests[c.Challenge] = *c
	return nil
}

func (m *MemoryManager) GetConsentRequest(challenge string) (*ConsentRequest, error) {
	if c, ok := m.consentRequests[challenge]; ok {
		return &c, nil
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) HandleConsentRequest(challenge string, r *HandledConsentRequest) (*ConsentRequest, error) {
	m.handledConsentRequests[r.Challenge] = *r
	return m.GetConsentRequest(challenge)
}

func (m *MemoryManager) VerifyAndInvalidateConsentRequest(verifier string) (*HandledConsentRequest, error) {
	for _, c := range m.consentRequests {
		if c.Verifier == verifier {
			for _, h := range m.handledConsentRequests {
				if h.Challenge == c.Challenge {
					if h.WasUsed {
						return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already"))
					}

					h.WasUsed = true
					if _, err := m.HandleConsentRequest(h.Challenge, &h); err != nil {
						return nil, err
					}

					h.ConsentRequest = &c
					return &h, nil
				}
			}
		}
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) FindPreviouslyGrantedConsentRequests(client string, subject string) ([]HandledConsentRequest, error) {
	var rs []HandledConsentRequest
	for _, c := range m.handledConsentRequests {
		cr, err := m.GetConsentRequest(c.Challenge)
		if err != nil {
			return nil, err
		}

		if client != cr.Client.GetID() {
			continue
		}

		if subject != cr.Subject {
			continue
		}

		if c.Error != nil {
			continue
		}

		if !c.Remember {
			continue
		}

		if cr.Skip {
			continue
		}

		if c.RememberFor > 0 && c.RequestedAt.Add(time.Duration(c.RememberFor)*time.Second).Before(time.Now().UTC()) {
			continue
		}

		c.ConsentRequest = cr
		rs = append(rs, c)
	}
	if len(rs) == 0 {
		return []HandledConsentRequest{}, nil
	}

	return rs, nil
}

func (m *MemoryManager) GetAuthenticationSession(id string) (*AuthenticationSession, error) {
	if c, ok := m.authSessions[id]; ok {
		return &c, nil
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) CreateAuthenticationSession(a *AuthenticationSession) error {
	m.authSessions[a.ID] = *a
	return nil
}

func (m *MemoryManager) DeleteAuthenticationSession(id string) error {
	delete(m.authSessions, id)
	return nil
}

func (m *MemoryManager) CreateAuthenticationRequest(a *AuthenticationRequest) error {
	m.authRequests[a.Challenge] = *a
	return nil
}

func (m *MemoryManager) GetAuthenticationRequest(challenge string) (*AuthenticationRequest, error) {
	if c, ok := m.authRequests[challenge]; ok {
		return &c, nil
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) HandleAuthenticationRequest(challenge string, r *HandledAuthenticationRequest) (*AuthenticationRequest, error) {
	m.handledAuthRequests[r.Challenge] = *r
	return m.GetAuthenticationRequest(challenge)
}

func (m *MemoryManager) VerifyAndInvalidateAuthenticationRequest(verifier string) (*HandledAuthenticationRequest, error) {
	for _, c := range m.authRequests {
		if c.Verifier == verifier {
			for _, h := range m.handledAuthRequests {
				if h.Challenge == c.Challenge {
					if h.WasUsed {
						return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Authentication verifier has been used already"))
					}

					h.WasUsed = true
					if _, err := m.HandleAuthenticationRequest(h.Challenge, &h); err != nil {
						return nil, err
					}

					h.AuthenticationRequest = &c
					return &h, nil
				}
			}
		}
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}
