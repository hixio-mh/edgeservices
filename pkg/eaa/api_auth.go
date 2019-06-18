// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eaa

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/smartedgemec/appliance-ce/pkg/eva/internal_pb"
	"google.golang.org/grpc"
)

func validateAppIP(ipAddress string, validationEndpoint string) (bool, error) {

	// Dial to EVA to get Edge Application ID and use it for validation
	conn, err := grpc.Dial(validationEndpoint, grpc.WithInsecure())
	if err != nil {
		return false, errors.Wrapf(err,
			"Failed to create a connection to %s", validationEndpoint)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(),
		3*time.Second)
	defer cancel()

	client := pb.NewIPApplicationLookupServiceClient(conn)
	requestBody := pb.IPApplicationLookupInfo{
		IpAddress: ipAddress,
	}

	lookupResult, err := client.GetApplicationByIP(ctx, &requestBody,
		grpc.WaitForReady(true))
	if err != nil {
		return false, errors.Wrap(err, "Cannot get App ID from EVA")
	}

	return lookupResult.AppID != "", nil
}

// RequestCredentials handles PKI for an application
func RequestCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		identity    AuthIdentity
		credentials AuthCredentials
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewDecoder(r.Body).Decode(&identity)
	if err != nil {
		log.Errf("/Auth RequestCredentials decode failed: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Errf("Cannot retrieve IP from RemoteAddr: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isIPValid, err := validateAppIP(host, cfg.InternalEndpoint)
	if err != nil {
		log.Errf("IP address validation failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isIPValid {
		log.Info("IP address invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cert, err := SignCSR(identity.Csr)
	if err != nil {
		log.Errf("/Auth RequestCredentials failed: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signedCertBlock := pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if signedCertBlock == nil {
		log.Err("/Auth RequestCredentials failed to enode signed cert")
		return
	}
	rcaBlock := pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE",
			Bytes: eaaCtx.certsEaaCa.rca.x509Cert.Raw})
	if rcaBlock == nil {
		log.Err("/Auth RequestCredentials failed to enode rca cert")
		return
	}

	credentials.ID = cert.Subject.CommonName
	credentials.Certificate = string(signedCertBlock)
	credentials.CaChain = []string{string(rcaBlock)}
	credentials.CaPool = []string{string(rcaBlock)}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(credentials)
	if err != nil {
		log.Errf("/Auth RequestCredentials encoding output to JSON failed: %s",
			err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}