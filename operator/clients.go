package operator

import (
	"errors"

	"github.com/jnnkrdb/corerdb/prtcl"
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

	prtcl.Log.Println("initializing kubernetes api-connection")

	if _config, err = rest.InClusterConfig(); err != nil {

		prtcl.Log.Println("error while initialization:", err)

		prtcl.PrintObject(_config, err)

	} else {

		if _k8sclient, err = kubernetes.NewForConfig(_config); err != nil {

			prtcl.Log.Println("clientset error:", err)

			prtcl.PrintObject(_config, _k8sclient, err)

		} else {

			prtcl.Log.Println("created kubernetes.Clientset and rest.Config")
		}
	}

	return
}

// initialize the custom resource definition rest client, based on the generated kubernetes client
func InitCRDOperatorRestClient(groupname, groupversion string, schemeAddFunc func(s *runtime.Scheme) error) (err error) {

	prtcl.Log.Println("initializing crds rest api-connection")

	if _config == nil {

		err = errors.New("can not initialize crd restclient before kubernetes clientset")

		prtcl.Log.Println("error initializing crd-rest-client:", err)

	} else {

		if err = schemeAddFunc(scheme.Scheme); err != nil {

			prtcl.Log.Println("error adding scheme:", err)

		} else {

			// create the rest client for the crds
			crdConf := *_config

			crdConf.ContentConfig.GroupVersion = &schema.GroupVersion{
				Group:   groupname,
				Version: groupversion,
			}

			crdConf.APIPath = "/apis"

			crdConf.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)

			crdConf.UserAgent = rest.DefaultKubernetesUserAgent()

			if _crdclient, err = rest.UnversionedRESTClientFor(&crdConf); err != nil {

				prtcl.Log.Println("failed loading restclient:", err)

				prtcl.PrintObject(crdConf, _crdclient, err)
			}
		}
	}

	return
}
