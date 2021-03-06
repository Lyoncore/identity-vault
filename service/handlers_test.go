// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/snapcore/snapd/asserts"
)

func TestSignHandlerNilData(t *testing.T) {
	sendRequestSignError(t, "POST", "/1.0/sign", nil)
}

func TestSignHandlerNoData(t *testing.T) {
	sendRequestSignError(t, "POST", "/1.0/sign", new(bytes.Buffer))
}

func TestSignHandlerInactive(t *testing.T) {
	// Mock the database
	config := ConfigSettings{KeyStoreType: "filesystem", KeyStorePath: "../keystore"}
	Environ = &Env{DB: &mockDB{}, Config: config}
	Environ.KeypairDB, _ = GetKeyStore(config)

	const assertions = `type: serial
authority-id: System
brand-id: System
model: Inactive
serial: A1234/L
revision: 1
timestamp: 2016-01-02T15:04:05Z
device-key: openpgp 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

openpgp env.KeypairDB, err = service.GetKeyStore(env.Config)
`
	result, _ := sendRequestSignError(t, "POST", "/1.0/sign", bytes.NewBufferString(assertions))

	if result.ErrorCode != "error-model-not-active" {
		t.Errorf("Expected 'error-model-not-active', got %v", result.ErrorCode)
	}
}

func TestSignHandler(t *testing.T) {
	// Mock the database
	config := ConfigSettings{KeyStoreType: "filesystem", KeyStorePath: "../keystore"}
	Environ = &Env{DB: &mockDB{}, Config: config}
	Environ.KeypairDB, _ = GetKeyStore(config)

	const assertions = `type: serial
authority-id: System
brand-id: System
model: Alder
serial: A1234/L
revision: 1
timestamp: 2016-01-02T15:04:05Z
device-key: openpgp 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

openpgp env.KeypairDB, err = service.GetKeyStore(env.Config)
`
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/1.0/sign", bytes.NewBufferString(assertions))
	http.HandlerFunc(SignHandler).ServeHTTP(w, r)

	// Check that we have a assertion as a response
	if w.Code != http.StatusOK {
		t.Errorf("Expected success HTTP status, got: %d", w.Code)
	}
	if w.Header().Get("Content-Type") != asserts.MediaType {
		t.Errorf("Expected content-type %s, got: %s", asserts.MediaType, w.Header().Get("Content-Type"))
	}
}

func TestSignHandlerBadAssertion(t *testing.T) {
	// Mock the database
	config := ConfigSettings{KeyStoreType: "filesystem", KeyStorePath: "../keystore"}
	Environ = &Env{DB: &mockDB{}, Config: config}
	Environ.KeypairDB, _ = GetKeyStore(config)

	const assertions = `type: serial
authority-id: System
brand-id: Vendor
model: Alder
serial: A1234/L
revision: This should be numeric
timestamp: 2016-01-02T15:04:05Z
device-key: openpgp mQINBFaiIK4BEADHpUmhX1koBIprWkUDQbqFCKZBPvKbwRkU3v5LNmFZJYsjAV3TqhFBUp61AHpr5pvTMw3fJ8j3hoH1of+rq8DtPtijUpoEXLhprO1S8OYzMQZpXAm8NIFQEWvjJQIkS0tcDDl8yRIMa81QVFpwuJ8B8ZTmYscmXtZdjZ7tP5WMk+hJTecBmO8Z3ZhCdDV819DRf7O5BUMau2YkkXfHQIzwsvRcXhQJMFjItkrZi9IquuTaqYhRWvc9ehj58f0GzkBkABn3UYiu3SpzS6tp1fEjqSrzPLxtWXwZNrMSaQET1juycCpYlYZe30ri07uH7heCmu9/bt112nrxdLYodPevzqoL/WL2ZMYxsdYnk0p382gmdrCNzWqja2dVXLD4YrAyG6Sm+a256OG2Tf3l01zMZnazDbI8c5FQdTKr+w8ugBbJYtAUcvczFCqrLGDFY2dFiFyzrCZYR/ac0WWWWV3pjNLsi35wD4jTiPmHzkMY7r6SefUntfha45EPHeefdsRAqKS/i67XEUliTo3XgH+h8yhQLNs+2CQ2mZXQ2aAV6iDH4jnJG4XQQXlT4t8y4AT5E6hgcfCIEd5K22th7B26ee0PJ5FRzcJPCy9+rbMBE5uvkd7nPiV1IBK7PFvMQRdV3pQRE837N4kbJy0ohgSq+lI0267gWzwK2nrJqv0q5wARAQABtCdEYXMgS2V5IDxqYW1lcy5qZXN1ZGFzb25AY2Fub25pY2FsLmNvbT6JAjgEEwECACIFAlaiIK4CGwMGCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEGGr9YjlK+ejdZ4QAK/DuiaZxUDx2rvakOYdr8949AyKTYyKIr+ruDaliVIn3xqUPWPPCVAScuy4oK9nigj99lUC02WBclUZPtUOjAOWQKlWm1+liwdYfb7Q+iBo92FTBMiJdAt30hCkX8yzqOjSD0Qdi9Q0Qnmk3JFGPPpqq7oUsdaBM8tbnG92nsDzaibKG9QzSyt5+CfapxTVa1xScDf+kJ2cO6lsTFUfOu8LKUDPojdwExF1iOMDMK3II4S47I+OlDL3kbznFLYlxzYRGGmGUwjl/Q19HscvmfjfZSHUK4bZCeZFvJPmG+1mByk91CJtOZDmyW5+MNRpfA7fa6kCKkFssCEvJVPMUrHvV5xSGXMcAkFoKlGALMVRrpW6d0/rImlMc5chDODYOephpvUimHFEoqvvjziNuyTqpLsfpInvyviQ6W7LRoJd6iCDZTGXA2c630QYggM7ti4SQ6Db9kScqKtf1pKky0FGa7RHlFM1zAoz51dLng/a3P/fEuZW4fArS/KJoR0wuYyQHZuxRlUi4P3OhUA+3NDAP8cjYvcVzQw4ksCbqzVS9kQNfXqT5Feg0UAxXqg80bDdJhxCG0ZjeMOZNXqPNKLkjARMsr6NNenjtddmKuEyzg3jUg2TAS0fqIuPSR6V2ynGA9tMh+ImluHPU+N8+TMl9jBkITU8SojgHkytjFbcuQINBFaiIK4BEAC2KyWyIorcnFuuPSenOhwVacqHxLEfRoZ5lG3oHcEpE/3Cy6c+etYR3j7Vb724FxEV+bUQGOewb2bRxnx8pot2yoV9Q6pA6Mzr5mdVqo7cfTua3ijj4bZhxtEQ4qz2qBC3zsT151cDzcYSfaJT6uwhcmqLmDhjarfrSElSHYRx2IFYhEMKLz9rvVKCfYD/cHgjzeUDGGMHUcS95jrOQ4EaH0Ok3jKVyjwgR3/4F1iwZuGXTnJ0SY2mUHgQxcoBM7e1qoOC+l4dia3GMWOQVCqFhtWH+1W58JkrUZ5dqRtJ5hYREE5wzrl6I8GQhLc7lS477Z6dK47LAsc6SfAQjCzTpugF9QYssHrXfeC629ak13tbCTZLbKY0opE2QWJprbKCfHxtFeMvk/IgbnNsAVnKPBBpZMKApPdorBscILteywJJCtzefirNkLXEhdYd6BU83wLWtTxPXJ9w2hnPFBYlRDufetk9CveeyMPOUXgp9zF8qhSBdxZ4wSZKEbgvihD0faOP9P8qbq2sO4GzbahY5tSzac+Lb+JfcysckR6taGdW7TdmysJnmcUq+ZIdmMdQEH7rQvlFImZThpDVQbPWELqBkyrC9l8+0QZLmBK+VkYbgqTC7Euyl/ffMpAtRu3q5uUPEIdqXUijydOdMKt5NbBhuKrz1PdJG2XC+UPGxwARAQABiQIfBBgBAgAJBQJWoiCuAhsMAAoJEGGr9YjlK+ej3QYP/090qBvsjHpMguEA9roNjLoLlCbmYs/NSKB1WR/61CKD0dZjI0VHcL0uso9fo6FRN9HWMNbdlBVBM81D56UlAdD+u1hq4HtFF/knV0BceBGDL9W9Hne0ntoYYqHdB8QL4Wm84JVuK3CMvBYx3cUVhtwB7UsxdXd6ujmHDqm3yk439gwX5nbCzx1tMgLPywMQWP6n/qW/oGj6l0Smew4QQKWPjhy4JqB52irKxO/gRuAimYy3jW1ls0b4Lgfq1NT00HNGT/QrqYmqhDsYPfVDPxlEuVnbuc+V1YidCUbsdbkyTNmge/oyqKruxyQajG7faMquuNkrD9uxKbk5vEaiU91AomQo8TBUvklQ4p238pnJQMoM8eMlfB40GCNG0RY/X3w79/n2YgCQ8Y5N2wuPh9bw5xN1xnadliDnDz7G32nCHmdoTD7sfml8sUHmUZutu3D2KXXDj+WTS5SlXDAdnhIbmw5FbJnBCenNe4Xix5yAHOkz5ICdaLpv/297PmZT+tll3eFDXRWgMYGT8sHtdUrDsNry1d6pGDxuKXXeZMkrMkJxBuZUdYYLepsA2JPwDq5mgsCA89zKIjdhDdy3lXQGKXtBiOzOqApSmjlmCuqIg3w5/quLWmcKkh6mp2l1gSkAc3ImjHveEYdvpZpaQWk2yQ5xuSjIJvcEs1jwFtSj

openpgp env.KeypairDB, err = service.GetKeyStore(env.Config)
`
	sendRequestSignError(t, "POST", "/1.0/sign", bytes.NewBufferString(assertions))

}

func TestSignHandlerBadAssertionNoRevision(t *testing.T) {
	// Mock the database
	config := ConfigSettings{KeyStoreType: "filesystem", KeyStorePath: "../keystore"}
	Environ = &Env{DB: &mockDB{}, Config: config}
	Environ.KeypairDB, _ = GetKeyStore(config)

	const assertions = `type: serial
authority-id: System
brand-id: Vendor
model: Alder
serial: A1234/L
timestamp: 2016-01-02T15:04:05Z
device-key: openpgp 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

openpgp env.KeypairDB, err = service.GetKeyStore(env.Config)
`
	sendRequestSignError(t, "POST", "/1.0/sign", bytes.NewBufferString(assertions))

}

func TestSignHandlerBadAssertionWrongType(t *testing.T) {
	// Mock the database
	config := ConfigSettings{KeyStoreType: "filesystem", KeyStorePath: "../keystore"}
	Environ = &Env{DB: &mockDB{}, Config: config}
	Environ.KeypairDB, _ = GetKeyStore(config)

	const assertions = `type: model
authority-id: System
brand-id: System
model: Alder
serial: A1234/L
series: Alder
revision: 1
os: 14.04
architecture: i686
gadget: magic wand
kernel: 4.2.0-35-generic
store: Canonical
class: Class
allowed-modes: all
required-snaps: gadget
timestamp: 2016-01-02T15:04:05Z
device-key: openpgp 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

openpgp PvKbwRkU3v5LNmFZJYsjAV3TqhFBUp61AHpr5pvTMw3fJ8j3h
`
	result, _ := sendRequestSignError(t, "POST", "/1.0/sign", bytes.NewBufferString(assertions))

	if result.ErrorSubcode != "error-invalid-type" {
		t.Errorf("Expected an 'invalid type' message, got %s", result.ErrorSubcode)
	}
}

func TestSignHandlerNonExistentModel(t *testing.T) {
	// Mock the database, ot finding the model
	Environ = &Env{DB: &errorMockDB{}}

	const assertions = `type: serial
authority-id: System
brand-id: Vendor
model: Cannot Find This
serial: A1234/L
revision: 1
timestamp: 2016-01-02T15:04:05Z
device-key: openpgp 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

openpgp env.KeypairDB, err = service.GetKeyStore(env.Config)
`
	sendRequestSignError(t, "POST", "/1.0/sign", bytes.NewBufferString(assertions))
}

func TestSignHandlerErrorKeyStore(t *testing.T) {
	// Mock the database and the keystore
	config := ConfigSettings{KeyStoreType: "filesystem", KeyStorePath: "../keystore"}
	Environ = &Env{DB: &mockDB{}, Config: config}
	Environ.KeypairDB, _ = getErrorMockKeyStore(config)

	const assertions = `type: serial
authority-id: System
brand-id: System
model: Alder
serial: A1234/L
revision: 1
timestamp: 2016-01-02T15:04:05Z
device-key: openpgp mQINBFaiIK4BEADHpUmhX1koBIprWkUDQbqFCKZBPvKbwRkU3v5LNmFZJYsjAV3TqhFBUp61AHpr5pvTMw3fJ8j3hoH1of+rq8DtPtijUpoEXLhprO1S8OYzMQZpXAm8NIFQEWvjJQIkS0tcDDl8yRIMa81QVFpwuJ8B8ZTmYscmXtZdjZ7tP5WMk+hJTecBmO8Z3ZhCdDV819DRf7O5BUMau2YkkXfHQIzwsvRcXhQJMFjItkrZi9IquuTaqYhRWvc9ehj58f0GzkBkABn3UYiu3SpzS6tp1fEjqSrzPLxtWXwZNrMSaQET1juycCpYlYZe30ri07uH7heCmu9/bt112nrxdLYodPevzqoL/WL2ZMYxsdYnk0p382gmdrCNzWqja2dVXLD4YrAyG6Sm+a256OG2Tf3l01zMZnazDbI8c5FQdTKr+w8ugBbJYtAUcvczFCqrLGDFY2dFiFyzrCZYR/ac0WWWWV3pjNLsi35wD4jTiPmHzkMY7r6SefUntfha45EPHeefdsRAqKS/i67XEUliTo3XgH+h8yhQLNs+2CQ2mZXQ2aAV6iDH4jnJG4XQQXlT4t8y4AT5E6hgcfCIEd5K22th7B26ee0PJ5FRzcJPCy9+rbMBE5uvkd7nPiV1IBK7PFvMQRdV3pQRE837N4kbJy0ohgSq+lI0267gWzwK2nrJqv0q5wARAQABtCdEYXMgS2V5IDxqYW1lcy5qZXN1ZGFzb25AY2Fub25pY2FsLmNvbT6JAjgEEwECACIFAlaiIK4CGwMGCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEGGr9YjlK+ejdZ4QAK/DuiaZxUDx2rvakOYdr8949AyKTYyKIr+ruDaliVIn3xqUPWPPCVAScuy4oK9nigj99lUC02WBclUZPtUOjAOWQKlWm1+liwdYfb7Q+iBo92FTBMiJdAt30hCkX8yzqOjSD0Qdi9Q0Qnmk3JFGPPpqq7oUsdaBM8tbnG92nsDzaibKG9QzSyt5+CfapxTVa1xScDf+kJ2cO6lsTFUfOu8LKUDPojdwExF1iOMDMK3II4S47I+OlDL3kbznFLYlxzYRGGmGUwjl/Q19HscvmfjfZSHUK4bZCeZFvJPmG+1mByk91CJtOZDmyW5+MNRpfA7fa6kCKkFssCEvJVPMUrHvV5xSGXMcAkFoKlGALMVRrpW6d0/rImlMc5chDODYOephpvUimHFEoqvvjziNuyTqpLsfpInvyviQ6W7LRoJd6iCDZTGXA2c630QYggM7ti4SQ6Db9kScqKtf1pKky0FGa7RHlFM1zAoz51dLng/a3P/fEuZW4fArS/KJoR0wuYyQHZuxRlUi4P3OhUA+3NDAP8cjYvcVzQw4ksCbqzVS9kQNfXqT5Feg0UAxXqg80bDdJhxCG0ZjeMOZNXqPNKLkjARMsr6NNenjtddmKuEyzg3jUg2TAS0fqIuPSR6V2ynGA9tMh+ImluHPU+N8+TMl9jBkITU8SojgHkytjFbcuQINBFaiIK4BEAC2KyWyIorcnFuuPSenOhwVacqHxLEfRoZ5lG3oHcEpE/3Cy6c+etYR3j7Vb724FxEV+bUQGOewb2bRxnx8pot2yoV9Q6pA6Mzr5mdVqo7cfTua3ijj4bZhxtEQ4qz2qBC3zsT151cDzcYSfaJT6uwhcmqLmDhjarfrSElSHYRx2IFYhEMKLz9rvVKCfYD/cHgjzeUDGGMHUcS95jrOQ4EaH0Ok3jKVyjwgR3/4F1iwZuGXTnJ0SY2mUHgQxcoBM7e1qoOC+l4dia3GMWOQVCqFhtWH+1W58JkrUZ5dqRtJ5hYREE5wzrl6I8GQhLc7lS477Z6dK47LAsc6SfAQjCzTpugF9QYssHrXfeC629ak13tbCTZLbKY0opE2QWJprbKCfHxtFeMvk/IgbnNsAVnKPBBpZMKApPdorBscILteywJJCtzefirNkLXEhdYd6BU83wLWtTxPXJ9w2hnPFBYlRDufetk9CveeyMPOUXgp9zF8qhSBdxZ4wSZKEbgvihD0faOP9P8qbq2sO4GzbahY5tSzac+Lb+JfcysckR6taGdW7TdmysJnmcUq+ZIdmMdQEH7rQvlFImZThpDVQbPWELqBkyrC9l8+0QZLmBK+VkYbgqTC7Euyl/ffMpAtRu3q5uUPEIdqXUijydOdMKt5NbBhuKrz1PdJG2XC+UPGxwARAQABiQIfBBgBAgAJBQJWoiCuAhsMAAoJEGGr9YjlK+ej3QYP/090qBvsjHpMguEA9roNjLoLlCbmYs/NSKB1WR/61CKD0dZjI0VHcL0uso9fo6FRN9HWMNbdlBVBM81D56UlAdD+u1hq4HtFF/knV0BceBGDL9W9Hne0ntoYYqHdB8QL4Wm84JVuK3CMvBYx3cUVhtwB7UsxdXd6ujmHDqm3yk439gwX5nbCzx1tMgLPywMQWP6n/qW/oGj6l0Smew4QQKWPjhy4JqB52irKxO/gRuAimYy3jW1ls0b4Lgfq1NT00HNGT/QrqYmqhDsYPfVDPxlEuVnbuc+V1YidCUbsdbkyTNmge/oyqKruxyQajG7faMquuNkrD9uxKbk5vEaiU91AomQo8TBUvklQ4p238pnJQMoM8eMlfB40GCNG0RY/X3w79/n2YgCQ8Y5N2wuPh9bw5xN1xnadliDnDz7G32nCHmdoTD7sfml8sUHmUZutu3D2KXXDj+WTS5SlXDAdnhIbmw5FbJnBCenNe4Xix5yAHOkz5ICdaLpv/297PmZT+tll3eFDXRWgMYGT8sHtdUrDsNry1d6pGDxuKXXeZMkrMkJxBuZUdYYLepsA2JPwDq5mgsCA89zKIjdhDdy3lXQGKXtBiOzOqApSmjlmCuqIg3w5/quLWmcKkh6mp2l1gSkAc3ImjHveEYdvpZpaQWk2yQ5xuSjIJvcEs1jwFtSj

openpgp env.KeypairDB, err = service.GetKeyStore(env.Config)
`
	result, _ := sendRequestSignError(t, "POST", "/1.0/sign", bytes.NewBufferString(assertions))

	if result.ErrorCode != "error-signing-assertions" {
		t.Errorf("Expected an 'error signing' message, got %s", result.ErrorCode)
	}
}

func TestVersionHandler(t *testing.T) {

	config := ConfigSettings{Version: "1.2.5"}
	Environ = &Env{Config: config}

	result, _ := sendRequestVersion(t, "GET", "/1.0/version", nil)

	if result.Version != Environ.Config.Version {
		t.Errorf("Incorrect version returned. Expected '%s' got: %v", Environ.Config.Version, result.Version)
	}

}

func TestModelsHandler(t *testing.T) {

	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/1.0/models", nil)
	http.HandlerFunc(ModelsHandler).ServeHTTP(w, r)

	// Check the JSON response
	result := ModelsResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the models response: %v", err)
	}
	if len(result.Models) != 6 {
		t.Errorf("Expected 6 models, got %d", len(result.Models))
	}
	if result.Models[0].Name != "Alder" {
		t.Errorf("Expected model name 'Alder', got %s", result.Models[0].Name)
	}
}

func TestModelsHandlerWithError(t *testing.T) {

	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/1.0/models", nil)
	http.HandlerFunc(ModelsHandler).ServeHTTP(w, r)

	// Check the JSON response
	result := ModelsResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the models response: %v", err)
	}
	if result.Success {
		t.Error("Expected error, got success response")
	}

}

func TestModelGetHandler(t *testing.T) {

	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	result, _ := sendRequest(t, "GET", "/1.0/models/1", nil)

	if result.Model.ID != 1 {
		t.Errorf("Expected model with ID 1, got %d", result.Model.ID)
	}
	if result.Model.Name != "Alder" {
		t.Errorf("Expected model name 'Alder', got %s", result.Model.Name)
	}
}

func TestModelGetHandlerWithError(t *testing.T) {

	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	sendRequestExpectError(t, "GET", "/1.0/models/999999", nil)
}

func TestModelGetHandlerWithBadID(t *testing.T) {

	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	sendRequestExpectError(t, "GET", "/1.0/models/999999999999999999999999999999", nil)
}

func TestModelUpdateHandler(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	// Update a model
	data := `
	{
	  "id": 1,
	  "brand-id": "System",
    "model":"聖誕快樂",
    "serial":"A1234/L",
		"revision": 2,
    "device-key":"ssh-rsa NNhqloxPyIYXiTP+3JTPWV/mNoBar2geWIf"
  }`

	result, _ := sendRequest(t, "PUT", "/1.0/models/1", bytes.NewBufferString(data))

	if result.Model.ID != 1 {
		t.Errorf("Expected model with ID 1, got %d", result.Model.ID)
	}
	if result.Model.Name != "聖誕快樂" {
		t.Errorf("Expected model name '聖誕快樂', got %s", result.Model.Name)
	}
}

func TestModelUpdateHandlerWithErrors(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	// Update a model
	data := `{}`

	sendRequestExpectError(t, "PUT", "/1.0/models/1", bytes.NewBufferString(data))
}

func TestModelUpdateHandlerWithNilData(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendRequestExpectError(t, "PUT", "/1.0/models/1", nil)
}

func TestModelUpdateHandlerWithEmptyData(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendRequestExpectError(t, "PUT", "/1.0/models/1", bytes.NewBufferString(""))
}

func TestModelUpdateHandlerWithBadData(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendRequestExpectError(t, "PUT", "/1.0/models/1", bytes.NewBufferString("bad"))
}

func TestModelUpdateHandlerWithBadID(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendRequestExpectError(t, "PUT", "/1.0/models/999999999999999999999999999999", bytes.NewBufferString("bad"))
}

func TestModelDeleteHandler(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	// Delete a model
	data := "{}"
	sendRequest(t, "DELETE", "/1.0/models/1", bytes.NewBufferString(data))
}

func TestModelDeleteHandlerWithErrors(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	// Delete a model
	data := `{}`

	sendRequestExpectError(t, "DELETE", "/1.0/models/1", bytes.NewBufferString(data))
}

func TestModelDeleteHandlerWithBadID(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendRequestExpectError(t, "DELETE", "/1.0/models/999999999999999999999999999999", bytes.NewBufferString("bad"))
}

func TestModelCreateHandler(t *testing.T) {
	// Mock the database
	config := ConfigSettings{KeyStoreType: "filesystem", KeyStorePath: "../keystore"}
	Environ = &Env{DB: &mockDB{}, Config: config}

	// Define a model linked with the signing-key as JSON
	model := ModelSerialize{BrandID: "System", Name: "聖誕快樂", Revision: 2, KeypairID: 1}
	data, _ := json.Marshal(model)

	result, _ := sendRequest(t, "POST", "/1.0/models", bytes.NewReader(data))
	if result.Model.ID != 7 {
		t.Errorf("Expected model with ID 7, got %d", result.Model.ID)
	}
	if result.Model.Name != "聖誕快樂" {
		t.Errorf("Expected model name '聖誕快樂', got %s", result.Model.Name)
	}
}

func TestModelCreateHandlerWithError(t *testing.T) {
	// Mock the database
	config := ConfigSettings{KeyStoreType: "filesystem", KeyStorePath: "../keystore"}
	Environ = &Env{DB: &errorMockDB{}, Config: config}

	// Define a model linked with the signing-key as JSON
	model := ModelSerialize{BrandID: "System", Name: "聖誕快樂", Revision: 2, KeypairID: 1}
	data, _ := json.Marshal(model)

	sendRequestExpectError(t, "POST", "/1.0/models", bytes.NewReader(data))
}

func TestModelCreateHandlerWithBase64Error(t *testing.T) {
	// Mock the database
	config := ConfigSettings{KeyStoreType: "filesystem"}
	Environ = &Env{DB: &errorMockDB{}, Config: config}

	// Define a model linked with the signing-key as JSON
	model := ModelSerialize{BrandID: "System", Name: "聖誕快樂", Revision: 2, KeypairID: 1}
	data, _ := json.Marshal(model)

	sendRequestExpectError(t, "POST", "/1.0/models", bytes.NewReader(data))
}

func TestModelCreateHandlerWithNilData(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendRequestExpectError(t, "POST", "/1.0/models", nil)
}

func TestModelCreateHandlerWithEmptyData(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendRequestExpectError(t, "POST", "/1.0/models", bytes.NewBufferString(""))
}

func TestModelCreateHandlerWithBadData(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendRequestExpectError(t, "POST", "/1.0/models", bytes.NewBufferString("bad"))
}

func sendRequest(t *testing.T, method, url string, data io.Reader) (ModelResponse, error) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, data)
	AdminRouter(Environ).ServeHTTP(w, r)

	// Check the JSON response
	result := ModelResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the model response: %v", err)
	}
	if !result.Success {
		t.Errorf("Expected success, got error: %s", result.ErrorMessage)
	}

	return result, err
}

func sendRequestExpectError(t *testing.T, method, url string, data io.Reader) (ModelResponse, error) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, data)
	AdminRouter(Environ).ServeHTTP(w, r)

	// Check the JSON response
	result := ModelResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the model response: %v", err)
	}
	if result.Success {
		t.Error("Expected error, got success")
	}

	return result, err
}

func sendRequestVersion(t *testing.T, method, url string, data io.Reader) (VersionResponse, error) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, data)
	SigningRouter(Environ).ServeHTTP(w, r)

	// Check the JSON response
	result := VersionResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the version response: %v", err)
	}

	return result, err
}

func sendRequestSignError(t *testing.T, method, url string, data io.Reader) (SignResponse, error) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, data)
	SigningRouter(Environ).ServeHTTP(w, r)

	if w.Code == http.StatusOK {
		t.Errorf("Expected error HTTP status, got: %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json; charset=UTF-8" {
		t.Errorf("Expected JSON content-type, got: %s", w.Header().Get("Content-Type"))
	}

	// Check the JSON response
	result := SignResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the signed response: %v", err)
	}
	if result.Success {
		t.Error("Expected an error, got success response")
	}

	return result, err
}
