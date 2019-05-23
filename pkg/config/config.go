/*
Copyright 2018 Heptio Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"path"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/heptio/sonobuoy/pkg/buildinfo"
	"github.com/heptio/sonobuoy/pkg/plugin"
	uuid "github.com/satori/go.uuid"
)

const (
	// DefaultNamespace is the namespace where the master and plugin workers will run (but not necessarily the pods created by the plugin workers).
	DefaultNamespace = "heptio-sonobuoy"

	// DefaultKubeConformanceImageURL is the URL of the docker image to run for the kube conformance tests.
	DefaultKubeConformanceImageURL = "gcr.io/heptio-images/kube-conformance"
	// UpstreamKubeConformanceImageURL is the URL of the docker image to run for
	// the kube conformance tests which is maintained by upstream Kubernetes.
	UpstreamKubeConformanceImageURL = "gcr.io/google-containers/conformance"
	// DefaultKubeConformanceImageTag is the default tag of the conformance image
	DefaultKubeConformanceImageTag = "latest"
	// DefaultAggregationServerBindPort is the default port for the aggregation server to bind to.
	DefaultAggregationServerBindPort = 8080
	// DefaultAggregationServerBindAddress is the default address for the aggregation server to bind to.
	DefaultAggregationServerBindAddress = "0.0.0.0"
	// DefaultAggregationServerTimeoutSeconds is the default amount of time the aggregation server will wait for all plugins to complete.
	DefaultAggregationServerTimeoutSeconds = 10800 // 180 min
	// MasterPodName is the name of the main pod that runs plugins and collects results.
	MasterPodName = "sonobuoy"
	// MasterContainerName is the name of the main container in the master pod.
	MasterContainerName = "kube-sonobuoy"
	// MasterResultsPath is the location in the main container of the master pod where results will be archived.
	MasterResultsPath = "/tmp/sonobuoy"
	// DefaultSonobuoyPullPolicy is the default pull policy used in the Sonobuoy config.
	DefaultSonobuoyPullPolicy = "IfNotPresent"
	// DefaultQueryQPS is the number of queries per second Sonobuoy will make when gathering data.
	DefaultQueryQPS = 30
	// DefaultQueryBurst is the peak number of queries per second Sonobuoy will make when gathering data.
	DefaultQueryBurst = 50
)

var (
	// DefaultKubeConformanceImage is the URL and tag of the docker image to run for the kube conformance tests.
	DefaultKubeConformanceImage = DefaultKubeConformanceImageURL + ":" + DefaultKubeConformanceImageTag
	// DefaultImage is the URL of the docker image to run for the aggregator and workers
	DefaultImage = "gcr.io/heptio-images/sonobuoy:" + buildinfo.Version
	// DefaultResources is the default set of resources which are queried for after plugins run. The strings
	// are compared against the resource.Name given by the client-go discovery client. The non-standard values
	// that are included here are: podlogs, servergroups, serverversion. The value 'nodes', although a crawlable
	// API value, also is used to query against the healthz and configz endpoints on the node.
	DefaultResources = []string{
		"apiservices",
		"certificatesigningrequests",
		"clusterrolebindings",
		"clusterroles",
		"componentstatuses",
		"configmaps",
		"controllerrevisions",
		"cronjobs",
		"customresourcedefinitions",
		"daemonsets",
		"deployments",
		"endpoints",
		"ingresses",
		"jobs",
		"leases",
		"limitranges",
		"mutatingwebhookconfigurations",
		"namespaces",
		"networkpolicies",
		"nodes",
		"persistentvolumeclaims",
		"persistentvolumes",
		"poddisruptionbudgets",
		"pods",
		"podsecuritypolicies",
		"podtemplates",
		"priorityclasses",
		"replicasets",
		"replicationcontrollers",
		"resourcequotas",
		"rolebindings",
		"roles",
		"servergroups",
		"serverversion",
		"serviceaccounts",
		"services",
		"statefulsets",
		"storageclasses",
		"validatingwebhookconfigurations",
		"volumeattachments",
	}
)

// FilterOptions allow operators to select sets to include in a report
type FilterOptions struct {
	Namespaces    string `json:"Namespaces"`
	LabelSelector string `json:"LabelSelector"`
}

// Config is the input struct used to determine what data to collect.
type Config struct {
	// NOTE: viper uses "mapstructure" as the tag for config
	// serialization, *NOT* "json".  mapstructure is a separate library
	// that converts maps to structs, and has its own syntax for tagging
	// fields. The only documentation on this is in the mapstructure docs:
	// https://godoc.org/github.com/mitchellh/mapstructure#example-Decode--Tags
	//
	// To be safe we annotate with both json and mapstructure tags.

	///////////////////////////////////////////////
	// Meta-Data collection options
	///////////////////////////////////////////////
	Description string `json:"Description" mapstructure:"Description"`
	UUID        string `json:"UUID" mapstructure:"UUID"`
	Version     string `json:"Version" mapstructure:"Version"`
	ResultsDir  string `json:"ResultsDir" mapstructure:"ResultsDir"`

	///////////////////////////////////////////////
	// Query options
	///////////////////////////////////////////////
	Resources []string      `json:"Resources" mapstructure:"Resources"`
	Filters   FilterOptions `json:"Filters" mapstructure:"Filters"`
	Limits    LimitConfig   `json:"Limits" mapstructure:"Limits"`
	QPS       float32       `json:"QPS,omitempty" mapstructure:"QPS"`
	Burst     int           `json:"Burst,omitempty" mapstructure:"Burst"`

	///////////////////////////////////////////////
	// Plugin configurations settings
	///////////////////////////////////////////////
	Aggregation      plugin.AggregationConfig `json:"Server" mapstructure:"Server"`
	PluginSelections []plugin.Selection       `json:"Plugins" mapstructure:"Plugins"`
	PluginSearchPath []string                 `json:"PluginSearchPath" mapstructure:"PluginSearchPath"`
	Namespace        string                   `json:"Namespace" mapstructure:"Namespace"`
	LoadedPlugins    []plugin.Interface       `json:"-"` // this is assigned when plugins are loaded.

	///////////////////////////////////////////////
	// Sonobuoy configuration
	///////////////////////////////////////////////
	WorkerImage       string            `json:"WorkerImage" mapstructure:"WorkerImage"`
	ImagePullPolicy   string            `json:"ImagePullPolicy" mapstructure:"ImagePullPolicy"`
	ImagePullSecrets  string            `json:"ImagePullSecrets" mapstructure:"ImagePullSecrets"`
	CustomAnnotations map[string]string `json:"CustomAnnotations,omitempty" mapstructure:"CustomAnnotations"`
}

// LimitConfig is a configuration on the limits of sizes of various responses.
type LimitConfig struct {
	PodLogs SizeOrTimeLimitConfig `json:"PodLogs" mapstructure:"PodLogs"`
}

// SizeOrTimeLimitConfig represents configuration that limits the size of
// something either by a total disk size, or by a length of time.
type SizeOrTimeLimitConfig struct {
	LimitSize string `json:"LimitSize" mapstructure:"LimitSize"`
	LimitTime string `json:"LimitTime" mapstructure:"LimitTime"`
}

// FilterResources is a utility function used to parse Resources
func (cfg *Config) FilterResources(filter []string) []string {
	var results []string
	for _, felement := range filter {
		for _, check := range cfg.Resources {
			if felement == check {
				results = append(results, felement)
			}
		}
	}
	return results
}

// OutputDir returns the directory under the ResultsDir containing the
// UUID for this run.
func (cfg *Config) OutputDir() string {
	return path.Join(cfg.ResultsDir, cfg.UUID)
}

// SizeLimitBytes returns how many bytes the configuration is set to limit,
// returning defaultVal if not set.
func (c SizeOrTimeLimitConfig) SizeLimitBytes(defaultVal int64) int64 {
	val, defaulted, err := c.sizeLimitBytes()

	// Ignore error, since we should have already caught it in validation
	if err != nil || defaulted {
		return defaultVal
	}

	return val
}

func (c SizeOrTimeLimitConfig) sizeLimitBytes() (val int64, defaulted bool, err error) {
	str := c.LimitSize
	if str == "" {
		return 0, true, nil
	}

	var bs datasize.ByteSize
	err = bs.UnmarshalText([]byte(str))
	return int64(bs.Bytes()), false, err
}

// TimeLimitDuration returns the duration the configuration is set to limit, returning defaultVal if not set.
func (c SizeOrTimeLimitConfig) TimeLimitDuration(defaultVal time.Duration) time.Duration {
	val, defaulted, err := c.timeLimitDuration()

	// Ignore error, since we should have already caught it in validation
	if err != nil || defaulted {
		return defaultVal
	}

	return val
}

func (c SizeOrTimeLimitConfig) timeLimitDuration() (val time.Duration, defaulted bool, err error) {
	str := c.LimitTime
	if str == "" {
		return 0, true, nil
	}

	val, err = time.ParseDuration(str)
	return val, false, err
}

// New returns a newly-constructed Config object with default values.
func New() *Config {
	var cfg Config
	cfgUuid, _ := uuid.NewV4()
	cfg.UUID = cfgUuid.String()
	cfg.Description = "DEFAULT"
	cfg.ResultsDir = "/tmp/sonobuoy"
	cfg.Version = buildinfo.Version

	cfg.Filters.Namespaces = ".*"

	cfg.QPS = DefaultQueryQPS
	cfg.Burst = DefaultQueryBurst
	cfg.Resources = DefaultResources

	cfg.Namespace = DefaultNamespace

	cfg.Aggregation.BindAddress = DefaultAggregationServerBindAddress
	cfg.Aggregation.BindPort = DefaultAggregationServerBindPort
	cfg.Aggregation.TimeoutSeconds = DefaultAggregationServerTimeoutSeconds

	cfg.PluginSearchPath = []string{
		"./plugins.d",
		"/etc/sonobuoy/plugins.d",
		"~/sonobuoy/plugins.d",
	}

	cfg.WorkerImage = DefaultImage
	cfg.ImagePullPolicy = DefaultSonobuoyPullPolicy

	return &cfg
}

// addPlugin adds a (configured, initialized) plugin to the config object so
// that it can be executed.
func (cfg *Config) addPlugin(plugin plugin.Interface) {
	cfg.LoadedPlugins = append(cfg.LoadedPlugins, plugin)
}

// getPlugins gets the list of plugins selected for this configuration.
func (cfg *Config) getPlugins() []plugin.Interface {
	return cfg.LoadedPlugins
}
