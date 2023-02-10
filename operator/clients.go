package operator

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var (
	_config *rest.Config = nil

	// kubernetes client
	_k8sclient *kubernetes.Clientset = nil

	// custom resource definition based client
	_crdclient *rest.RESTClient = nil
)

// get the current kubernetes client
func K8S() *kubernetes.Clientset {

	return _k8sclient
}

// get the current kubernetes client
func CRD() *rest.RESTClient {

	return _crdclient
}

// initialize the kubernetes client
func InitK8sOperatorClient() (err error) {

	if _config, err = rest.InClusterConfig(); err == nil {

		_k8sclient, err = kubernetes.NewForConfig(_config)
	}

	return
}

// initialize the custom resource definition rest client, based on the generated kubernetes client
//
// if no kuberntes client was generated before, it will be generated first
func InitCRDOperatorRestClient(groupname, groupversion string, schemeAddFunc func(s *runtime.Scheme) error) (err error) {

	if _config == nil {

		err = InitK8sOperatorClient()
	}

	if err == nil {

		if err = schemeAddFunc(scheme.Scheme); err == nil {

			// create the rest client for the crds
			crdConf := *_config

			crdConf.ContentConfig.GroupVersion = &schema.GroupVersion{
				Group:   groupname,
				Version: groupversion,
			}

			crdConf.APIPath = "/apis"

			crdConf.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)

			crdConf.UserAgent = rest.DefaultKubernetesUserAgent()

			_crdclient, err = rest.UnversionedRESTClientFor(&crdConf)
		}
	}

	return
}
