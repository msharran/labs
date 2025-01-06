"""A Kubernetes Python Pulumi program"""

import pulumi

import pulumi_kubernetes as kubernetes

certman = kubernetes.helm.v4.Chart("cert-manager",
                                   namespace="cert-manager",
                                   chart="oci://registry-1.docker.io/bitnamicharts/cert-manager",
                                   version="1.3.1")
