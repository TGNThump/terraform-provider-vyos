package vyos

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/foltik/vyos-client-go/client"
)

type VyosConfig struct {
	apiClient    *client.Client
	skipSaving   bool
	saveFile     string
	mutex        sync.Mutex
	cachedConfig *map[string]any
}

func New(apiClient *client.Client, skipSaving bool, saveFile string) *VyosConfig {
	config := &VyosConfig{
		skipSaving: skipSaving,
		saveFile:   saveFile,
		apiClient:  apiClient,
	}
	return config
}

func (vc *VyosConfig) SaveIfRequired(ctx context.Context) error {
	if vc.skipSaving {
		return nil
	} else if vc.saveFile == "" {
		return vc.apiClient.Config.Save(ctx)
	} else {
		return vc.apiClient.Config.SaveFile(ctx, vc.saveFile)
	}
}

func (vc *VyosConfig) getRemoteConfig(ctx context.Context) (map[string]any, error) {
	resp, err := vc.apiClient.Request(ctx, "retrieve", map[string]any{
		"op":   "showConfig",
		"path": []string{},
	})
	if err != nil {
		if strings.Contains(err.Error(), "could not fetch config") {
			// If we get an empty path error, consume it and return nil
			return nil, nil
		} else {
			return nil, err
		}
	}

	obj, ok := resp.(map[string]any)
	if !ok {
		return nil, errors.New("received unexpected repsonse format from server")
	}
	return obj, nil
}

func (vc *VyosConfig) GetFullConfig(ctx context.Context) (*map[string]any, error) {
	vc.mutex.Lock()
	if vc.cachedConfig == nil {
		config, err := vc.getRemoteConfig(ctx)
		if err != nil {
			vc.mutex.Unlock()
			return nil, err
		}
		vc.cachedConfig = &config
	}
	vc.mutex.Unlock()
	return vc.cachedConfig, nil
}

func (vc *VyosConfig) invalidateConfigCache() {
	vc.mutex.Lock()
	vc.cachedConfig = nil
	vc.mutex.Unlock()
}

func (vc *VyosConfig) Show(ctx context.Context, path string) (any, error) {
	path_components := strings.Split(path, " ")
	if path == "" {
		path_components = []string{}
	}

	fullConfig, err := vc.GetFullConfig(ctx)
	if err != nil {
		return nil, err
	}

	config, err := getConfigFromPath(*fullConfig, path_components)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (vc *VyosConfig) Set(ctx context.Context, path string, value any) error {
	vc.invalidateConfigCache()

	err := vc.apiClient.Config.Set(ctx, path, value)
	if err != nil {
		return err
	}

	err = vc.SaveIfRequired(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (vc *VyosConfig) Delete(ctx context.Context, path string) error {
	vc.invalidateConfigCache()
	err := vc.apiClient.Config.Delete(ctx, path)
	if err != nil {
		return err
	}

	err = vc.SaveIfRequired(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (vc *VyosConfig) ApiRequest(ctx context.Context, endpoint string, payload any) (any, error) {
	return vc.apiClient.Request(ctx, endpoint, payload)
}

func getConfigFromPath(configTree map[string]interface{}, path_components []string) (rval interface{}, err error) {
	var ok bool

	if len(path_components) == 0 {
		return configTree, nil
	}

	if rval, ok = configTree[path_components[0]]; !ok {
		return nil, nil
	} else if len(path_components) == 1 { // we've reached the final path component
		return rval, nil
	} else if configTree, ok = rval.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("malformed configTree at %#v", rval)
	} else { // 1+ more path components
		return getConfigFromPath(configTree, path_components[1:])
	}
}
