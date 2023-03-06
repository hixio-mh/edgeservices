// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

module kube-rsu

require (
	k8s.io/client-go v0.20.0-alpha.2 // indirect
	rsu v0.0.0
)

replace rsu v0.0.0 => ./cmd

go 1.16
