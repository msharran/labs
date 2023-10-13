package main

import (
	"example.com/cdk8s-poc/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewNginxChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}

	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	// define resources here

	labels := &map[string]*string{"app": jsii.String("nginx")}

	k8s.NewKubeDeployment(chart, jsii.String("nginx"), &k8s.KubeDeploymentProps{
		Metadata: &k8s.ObjectMeta{
			Labels:    labels,
			Name:      jsii.String("nginx"),
			Namespace: jsii.String("sharran"),
		},
		Spec: &k8s.DeploymentSpec{
			Selector: &k8s.LabelSelector{
				MatchLabels: labels,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: labels,
				},
				Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{
						{
							Name:  jsii.String("nginx"),
							Image: jsii.String("nginx:latest"),
						},
					}},
			},
		},
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewNginxChart(app, "cdk8s-poc", nil)
	app.Synth()
}
