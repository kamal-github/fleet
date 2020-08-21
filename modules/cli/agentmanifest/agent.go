package agentmanifest

import (
	"context"
	"crypto/tls"
	fmt "fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/rancher/fleet/modules/cli/agentconfig"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rancher/fleet/pkg/agent"

	"github.com/pkg/errors"
	"github.com/rancher/fleet/modules/cli/pkg/client"
	fleet "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	"github.com/rancher/fleet/pkg/config"
	fleetcontrollers "github.com/rancher/fleet/pkg/generated/controllers/fleet.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/yaml"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Options struct {
	CA       []byte
	Host     string
	NoCA     bool
	NoCheck  bool
	Labels   map[string]string
	ClientID string
}

func AgentToken(ctx context.Context, controllerNamespace, kubeConfigFile string, client *client.Client, tokenName string, opts *Options) ([]runtime.Object, error) {
	token, err := getToken(ctx, tokenName, client)
	if err != nil {
		return nil, err
	}

	kubeConfig, err := getKubeConfig(kubeConfigFile, client.Namespace, token, opts.Host, opts.CA, opts.NoCA)
	if err != nil {
		return nil, err
	}

	if !opts.NoCheck {
		if err := testKubeConfig(kubeConfig, opts.Host); err != nil {
			return nil, fmt.Errorf("failed to testing kubeconfig: %w", err)
		}
	}

	return objects(controllerNamespace, kubeConfig), nil
}

func insecurePing(host string) {
	// I do this to make k3s generate a new SAN if it needs to
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	defer client.CloseIdleConnections()

	resp, err := client.Get(host)
	if err == nil {
		resp.Body.Close()
	}
}

func testKubeConfig(kubeConfig, host string) error {
	insecurePing(host)

	restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		return fmt.Errorf("failed to test kubeconfig: %w", err)
	}
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to test build client from kubeconfig: %w", err)
	}
	_, err = client.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to test connection to %s: %w", host, err)
	}
	return nil
}

func AgentManifest(ctx context.Context, systemNamespace, controllerNamespace string, cg *client.Getter, output io.Writer, tokenName string, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}

	client, err := cg.Get()
	if err != nil {
		return err
	}

	objs, err := AgentToken(ctx, controllerNamespace, cg.Kubeconfig, client, tokenName, opts)
	if err != nil {
		return err
	}

	agentConfig, err := agentconfig.AgentConfig(ctx, controllerNamespace, cg, &agentconfig.Options{
		Labels:   opts.Labels,
		ClientID: opts.ClientID,
	})
	if err != nil {
		return err
	}

	objs = append(objs, agentConfig...)

	cfg, err := config.Lookup(ctx, systemNamespace, config.ManagerConfigName, client.Core.ConfigMap())
	if err != nil {
		return err
	}

	objs = append(objs, agent.Manifest(controllerNamespace, cfg.AgentImage, cfg.AgentImagePullPolicy)...)

	data, err := yaml.Export(objs...)
	if err != nil {
		return err
	}

	_, err = output.Write(data)
	return err
}

func checkHost(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return errors.Wrapf(err, "invalid host, override with --server-url")
	}
	if u.Hostname() == "localhost" || strings.HasPrefix(u.Hostname(), "127.") || u.Hostname() == "0.0.0.0" {
		return fmt.Errorf("invalid host %s in server URL, use --server-url to set a proper server URL for the kubernetes endpoint", u.Hostname())
	}
	return nil
}

func getKubeConfig(kubeConfig string, namespace, token, host string, ca []byte, noCA bool) (string, error) {
	cc := kubeconfig.GetNonInteractiveClientConfig(kubeConfig)
	cfg, err := cc.RawConfig()
	if err != nil {
		return "", err
	}

	host, doCheckHost, err := getHost(host, cfg)
	if err != nil {
		return "", err
	}

	if doCheckHost {
		if err := checkHost(host); err != nil {
			return "", err
		}
	}

	if noCA {
		ca = nil
	} else {
		ca, err = getCA(ca, cfg)
		if err != nil {
			return "", err
		}
	}

	cfg = clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster": {
				Server:                   host,
				CertificateAuthorityData: ca,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user": {
				Token: token,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"default": {
				Cluster:   "cluster",
				AuthInfo:  "user",
				Namespace: namespace,
			},
		},
		CurrentContext: "default",
	}

	data, err := clientcmd.Write(cfg)
	return string(data), err
}

func getCluster(cfg clientcmdapi.Config) (*clientcmdapi.Cluster, error) {
	ctx := cfg.Contexts[cfg.CurrentContext]
	if ctx == nil {
		return nil, fmt.Errorf("failed to find host for agent access, context not found")
	}

	cluster := cfg.Clusters[ctx.Cluster]
	if cluster == nil {
		return nil, fmt.Errorf("failed to find host for agent access, cluster not found")
	}

	return cluster, nil
}

func getHost(host string, cfg clientcmdapi.Config) (string, bool, error) {
	if host != "" {
		return host, false, nil
	}

	cluster, err := getCluster(cfg)
	if err != nil {
		return "", false, err
	}

	return cluster.Server, true, nil
}

func getCA(ca []byte, cfg clientcmdapi.Config) ([]byte, error) {
	if len(ca) > 0 {
		return ca, nil
	}

	cluster, err := getCluster(cfg)
	if err != nil {
		return nil, err
	}

	if len(cluster.CertificateAuthorityData) > 0 {
		return cluster.CertificateAuthorityData, nil
	}

	if cluster.CertificateAuthority != "" {
		return ioutil.ReadFile(cluster.CertificateAuthority)
	}

	return nil, nil
}

func getToken(ctx context.Context, tokenName string, client *client.Client) (string, error) {
	secretName, err := waitForSecretName(ctx, tokenName, client)
	if err != nil {
		return "", err
	}

	secret, err := client.Core.Secret().Get(client.Namespace, secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	token := secret.Data[coreV1.ServiceAccountTokenKey]
	if len(token) == 0 {
		return "", fmt.Errorf("failed to find token on secret %s/%s", client.Namespace, secretName)
	}

	return string(token), nil
}

func waitForSecretName(ctx context.Context, tokenName string, client *client.Client) (string, error) {
	watcher, err := startWatch(client.Namespace, client.Fleet.ClusterRegistrationToken())
	if err != nil {
		return "", err
	}
	defer func() {
		watcher.Stop()
		for range watcher.ResultChan() {
			// drain the channel
		}
	}()

	crt, err := client.Fleet.ClusterRegistrationToken().Get(client.Namespace, tokenName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to lookup token %s: %w", tokenName, err)
	}
	if crt.Status.SecretName != "" {
		return crt.Status.SecretName, nil
	}

	timeout := time.After(time.Minute)
	for {
		var event watch.Event
		select {
		case <-timeout:
			return "", fmt.Errorf("timeout getting credential for cluster group")
		case <-ctx.Done():
			return "", ctx.Err()
		case event = <-watcher.ResultChan():
		}

		if newCGT, ok := event.Object.(*fleet.ClusterRegistrationToken); ok {
			if newCGT.UID != crt.UID || newCGT.Status.SecretName == "" {
				continue
			}
			return newCGT.Status.SecretName, nil
		}
	}
}

func startWatch(namespace string, sa fleetcontrollers.ClusterRegistrationTokenClient) (watch.Interface, error) {
	secrets, err := sa.List(namespace, metav1.ListOptions{
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	return sa.Watch(namespace, metav1.ListOptions{ResourceVersion: secrets.ResourceVersion})
}
