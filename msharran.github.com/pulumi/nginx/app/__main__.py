"""A Kubernetes Python Pulumi program"""

import pulumi
from pulumi_kubernetes.apps.v1 import Deployment
from pulumi_kubernetes.core.v1 import Service, Namespace

config = pulumi.Config()

name = "nginx-app"
app_labels = { "service": "nginx-app"}

namespace = Namespace(
    name,
    metadata={ "name": name })

service = Service(
    name,
    metadata={ "name": name, "namespace": namespace.metadata["name"] },
    spec={
        "ports": [{ "port": 80, "target_port": 80 }],
        "selector": app_labels,
    })

deployment = Deployment(
    name,
    metadata={ "name": name, "namespace": namespace.metadata["name"] },
    spec={
        "selector": { "match_labels": app_labels },
        "replicas": config.require_int("replicas"),
        "template": {
            "metadata": { "labels": app_labels },
            "spec": { "containers": [{ "name": name, "image": "nginx" }] }
        },
    })

pulumi.export("name", deployment.metadata["name"])
