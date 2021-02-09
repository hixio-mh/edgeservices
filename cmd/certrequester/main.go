// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package main

import (
	"flag"
	"os"

	"github.com/open-ness/common/log"
	"github.com/open-ness/edgenode/pkg/certrequester"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errf("Failed to get cluster config %s\n", err.Error())
		os.Exit(1)
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		log.Errf("Failed to initialize clientset: %s\n", err.Error())
		os.Exit(1)
	}

	configPath := flag.String("cfg", "certrequest.json", "CSR config path")
	flag.Parse()

	err = certrequester.GetCertificate(clientset, *configPath)
	if err != nil {
		log.Errf("Failed to generate certificate: %s\n", err.Error())
		os.Exit(1)
	}
}
