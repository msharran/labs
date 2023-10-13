package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "net/http/pprof"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/metadata/metadatainformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/pager"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil) // initializing the flags
	defer klog.Flush()  // flushes all pending log I/O
	flag.Parse()        // parses the command-line flags

	klog.Info("now you can see me")

	sigchan, ctx, cancel := setupSignalHandling()
	defer func() {
		signal.Stop(sigchan)
		cancel()
	}()

	go http.ListenAndServe("localhost:8080", nil)

	for i := 0; i < 1; i++ {
		go watchViaListWatch(ctx)
		//go watchViaSharedInformers(ctx)
		//go watchViaMetaInformers(ctx)
		//go watchViaReflector(ctx)
	}

	<-ctx.Done()

}

func watchViaListWatch(ctx context.Context) {
	clientset := kubernetes.NewForConfigOrDie(getKubeConfig())

	lw := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "secrets", corev1.NamespaceAll, fields.Everything())

	processEvents := func() {
		w, err := lw.Watch(metav1.ListOptions{LabelSelector: "example.com/managed=true"})
		if err != nil {
			panic(err)
		}

		for event := range w.ResultChan() {
			switch event.Type {
			case watch.Added:
				s := event.Object.(*v1.Secret)
				fmt.Printf("Added: %s/%s\n", s.ObjectMeta.Namespace, s.ObjectMeta.Name)
				// fmt.Printf("###%+v\n", s)
			case watch.Modified:
				s := event.Object.(*v1.Secret)
				fmt.Printf("Modified: %s/%s modified\n", s.ObjectMeta.Namespace, s.ObjectMeta.Name)
			case watch.Deleted:
				s := event.Object.(*v1.Secret)
				fmt.Printf("Deleted: %s/%s deleted\n", s.ObjectMeta.Namespace, s.ObjectMeta.Name)
			default:
				fmt.Printf("Error: %+v\n", event)
			}
		}
	}

	// ss := listSecrets(ctx, clientset)
	// for i, s := range ss {
	// 	fmt.Printf("####: %v = %v/%v\n", i, s.Namespace, s.Name)
	// }
	// if true {
	// 	return
	// }

	go func() {
		wait.Until(processEvents, time.Second, ctx.Done())
		fmt.Printf("exiting watch\n")
	}()

	<-ctx.Done()
}

func watchViaReflector(ctx context.Context) {
	clientset := kubernetes.NewForConfigOrDie(getKubeConfig())

	lw := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "secrets", metav1.NamespaceAll, fields.Everything())

	r := cache.NewReflector(
		lw,
		&corev1.Secret{},
		&cache.FakeCustomStore{
			AddFunc: func(obj interface{}) error {
				s := obj.(*v1.Secret)
				fmt.Printf("Added: %s/%s\n", s.ObjectMeta.Namespace, s.ObjectMeta.Name)
				return nil
			},
			UpdateFunc: func(obj interface{}) error {
				s := obj.(*v1.Secret)
				fmt.Printf("Updated: %s/%s\n", s.ObjectMeta.Namespace, s.ObjectMeta.Name)
				return nil
			},
			DeleteFunc: func(obj interface{}) error {
				s := obj.(*v1.Secret)
				fmt.Printf("Deleted: %s/%s\n", s.ObjectMeta.Namespace, s.ObjectMeta.Name)
				return nil
			},

			ListFunc:     func() []interface{} { panic("foo") },
			ListKeysFunc: func() []string { panic("foo") },
			GetFunc:      func(obj interface{}) (item interface{}, exists bool, err error) { panic("foo") },
			GetByKeyFunc: func(key string) (item interface{}, exists bool, err error) { panic("foo") },
			ReplaceFunc: func(list []interface{}, resourceVersion string) error {
				for _, o := range list {
					s := o.(*corev1.Secret)
					fmt.Printf("Replace: %s/%s\n", s.ObjectMeta.Namespace, s.ObjectMeta.Name)
				}
				return nil
			},
			ResyncFunc: func() error { panic("foo") },
		},
		0,
	)

	go r.Run(ctx.Done())

	<-ctx.Done()
}

func watchViaMetaInformers(ctx context.Context) {
	metaclient := metadata.NewForConfigOrDie(getKubeConfig())

	informerFactory := metadatainformer.NewFilteredSharedInformerFactory(metaclient, 6*time.Hour, metav1.NamespaceAll,
		func(options *metav1.ListOptions) {
			options.LabelSelector = "example.com/managed=true"
		})
	secretsInformer := informerFactory.ForResource(corev1.SchemeGroupVersion.WithResource("secrets"))

	count := 0
	secretsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			s := obj.(*metav1.PartialObjectMetadata)
			count += 1
			fmt.Printf("AddFunc: %v = %v/%v\n", count, s.Namespace, s.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			s := newObj.(*metav1.PartialObjectMetadata)
			fmt.Printf("UpdateFunc: %v/%v\n", s.Namespace, s.Name)
		},
		DeleteFunc: func(obj interface{}) {
			s := obj.(*metav1.PartialObjectMetadata)
			fmt.Printf("DeleteFunc: %v/%v\n", s.Namespace, s.Name)
		},
	})

	informerFactory.Start(ctx.Done())
	fmt.Println("waiting for caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), secretsInformer.Informer().HasSynced) {
		panic(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	fmt.Println("caches synced")

	<-ctx.Done()
}

func watchViaSharedInformers(ctx context.Context) {
	clientset := kubernetes.NewForConfigOrDie(getKubeConfig())

	informerFactory := informers.NewSharedInformerFactoryWithOptions(clientset, 6*time.Hour,
		informers.WithTweakListOptions(
			func(options *metav1.ListOptions) {
				options.LabelSelector = "example.com/managed=true"
				//options.FieldSelector = "metadata.name,metadata.namespace" //fields.OneTermEqualSelector("metadata.name", nodeName).String()
			}))
	secretsInformer := informerFactory.Core().V1().Secrets()

	// secretsInformer.Informer().SetTransform(func(obj interface{}) (interface{}, error) {
	// 	if s, ok := obj.(*corev1.Secret); ok {
	// 		ss := &corev1.Secret{
	// 			ObjectMeta: metav1.ObjectMeta{
	// 				Name:      s.GetName(),
	// 				Namespace: s.GetNamespace(),
	// 				// Labels:      s.GetLabels(),
	// 				// Annotations: s.GetAnnotations(),
	// 			},
	// 			Type: s.Type,
	// 		}
	// 		return ss, nil
	// 	}
	// 	return obj, nil
	// })

	count := 0
	secretsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			s := obj.(*corev1.Secret)
			count += 1
			fmt.Printf("AddFunc: %v = %v/%v\n", count, s.Namespace, s.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			s := newObj.(*corev1.Secret)
			fmt.Printf("UpdateFunc: %v/%v\n", s.Namespace, s.Name)
		},
		DeleteFunc: func(obj interface{}) {
			s := obj.(*corev1.Secret)
			fmt.Printf("DeleteFunc: %v/%v\n", s.Namespace, s.Name)
		},
	})

	informerFactory.Start(ctx.Done())
	fmt.Println("waiting for caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), secretsInformer.Informer().HasSynced) {
		panic(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	fmt.Println("caches synced")

	<-ctx.Done()
}

func getKubeConfig() *rest.Config {
	home := homedir.HomeDir()
	kc := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kc)
	if err != nil {
		panic(err)
	}
	config.Burst = 1000
	config.QPS = 500

	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return &debugRoundTripper{delegatedRoundTripper: rt}
	})
	return config
}

func setupSignalHandling() (chan os.Signal, context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	return c, ctx, cancel
}

type debugRoundTripper struct {
	delegatedRoundTripper http.RoundTripper
}

func (rt *debugRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var respStatus int
	resp, err := rt.delegatedRoundTripper.RoundTrip(req)
	if resp != nil {
		respStatus = resp.StatusCode
	}

	fmt.Printf("####: req-url=%v req-headers=%v resp-status=%v error=%v\n", req.URL.String(), req.Header, respStatus, err)

	//debug.PrintStack()

	return resp, err
}

func listSecrets(ctx context.Context, clientset kubernetes.Interface) []*corev1.Secret {
	opts := metav1.ListOptions{LabelSelector: "example.com/managed=true"}

	lp := pager.New(func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		return clientset.CoreV1().Secrets(metav1.NamespaceAll).List(ctx, opts)
	})

	// result, isPaged, err := lp.List(ctx, opts)
	// if err != nil {
	// 	panic(err)
	// }

	var ss []*corev1.Secret
	err := lp.EachListItem(ctx, opts, func(obj runtime.Object) error {
		s := obj.(*corev1.Secret)
		ss = append(ss, s)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return ss
}
