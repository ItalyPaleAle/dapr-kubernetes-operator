package helm

import (
	"context"
	"time"

	"github.com/dapr-sandbox/dapr-kubernetes-operator/pkg/utils/mergemap"
	"github.com/rs/xid"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
)

type InstallOption func(*ReleaseOptions[action.Install])

func (h *Helm) Install(ctx context.Context, chart string, options ...InstallOption) (*release.Release, error) {
	client := action.NewInstall(h.config)
	client.ReleaseName = xid.New().String()
	client.CreateNamespace = true
	client.Devel = true
	client.IncludeCRDs = true
	client.Wait = true
	client.Namespace = xid.New().String()
	client.Timeout = 10 * time.Minute

	io := ReleaseOptions[action.Install]{
		Client: client,
		Values: make(map[string]interface{}),
	}

	for _, option := range options {
		option(&io)
	}

	cp, err := client.ChartPathOptions.LocateChart(chart, h.settings)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	return client.RunWithContext(
		ctx,
		chartRequested,
		mergemap.Merge(map[string]interface{}{}, io.Values),
	)
}
