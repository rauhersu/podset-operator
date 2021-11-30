STEPS
=====

Based on: https://learn.openshift.com/operatorframework/go-operator-podset/

$ operator-sdk version
operator-sdk version: "v1.10.1-ocp", commit: "58bf7213d22688d7d0cd4ccc7f4faf4ec5b33e3f", kubernetes version: "v1.21", go version: "go1.16.6", GOOS: "linux", GOARCH: "amd64"

$ operator-sdk init --domain=example.com --repo=github.com/redhat/podset-operator
$ operator-sdk create api --group=app --version=v1alpha1 --kind=Podset --resource --controller
$ git init .
$ git branch -M main
$ git remote add origin https://github.com/rauhersu/podset-operator.git
$ make generate
$ make manifests
$ oc new-project myproject

Local:

$ WATCH_NAMESPACE=myproject make install run
$ oc project myproject
$ oc apply -f config/samples/app_v1alpha1_podset.yaml
$ oc get podset
NAME            DESIRED   AVAILABLE
podset-sample   3         3

In cluster:

$ IMAGE_TAG_BASE=docker.io/rauherna/podset-operator make docker-build docker-push
$ oc project myproject
$ IMAGE_TAG_BASE=docker.io/rauherna/podset-operator make deploy
