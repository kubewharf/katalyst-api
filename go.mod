module github.com/kubewharf/katalyst-api

go 1.17

require (
	github.com/gogo/protobuf v1.3.2
	google.golang.org/grpc v1.51.0
	k8s.io/api v0.24.6
	k8s.io/apimachinery v0.24.6
	k8s.io/client-go v0.18.10
	k8s.io/code-generator v0.18.10
	k8s.io/klog/v2 v2.70.1
	k8s.io/kubelet v0.18.10
	k8s.io/kubernetes v1.18.10
	k8s.io/metrics v0.18.10
	sigs.k8s.io/yaml v1.2.0
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.3.1-0.20221206200815-1e63c2f08a10 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/component-base v0.24.6 // indirect
	k8s.io/gengo v0.0.0-20211129171323-c02415ce4185 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.18.10
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.10
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.10
	k8s.io/apiserver => k8s.io/apiserver v0.18.10
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.10
	k8s.io/client-go => k8s.io/client-go v0.18.10
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.10
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.10
	k8s.io/code-generator => k8s.io/code-generator v0.18.10
	k8s.io/component-base => k8s.io/component-base v0.18.10
	k8s.io/component-helpers => k8s.io/component-helpers v0.18.10
	k8s.io/controller-manager => k8s.io/controller-manager v0.18.10
	k8s.io/cri-api => k8s.io/cri-api v0.18.10
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.10
	k8s.io/klog => k8s.io/klog v1.0.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.10
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.10
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.10
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.10
	k8s.io/kubectl => k8s.io/kubectl v0.18.10
	k8s.io/kubelet => github.com/kubewharf/kubelet v1.24.6-kubewharf.8
	k8s.io/kubernetes => k8s.io/kubernetes v1.18.10
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.10
	k8s.io/metrics => k8s.io/metrics v0.18.10
	k8s.io/mount-utils => k8s.io/mount-utils v0.18.10
	k8s.io/node-api => k8s.io/node-api v0.18.10
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.18.10
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.10
)
