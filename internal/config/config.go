/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package config defines the config for cli.
package config

import (
	"encoding/json"
	"errors"
	"fmt"

	// for unmarshal yaml config file
	_ "gopkg.in/yaml.v2"

	"github.com/TencentBlueKing/bscp-go/internal/constant"
)

// ClientConfig config for bscp-go when run as daemon
type ClientConfig struct {
	ConfigFile string `json:"config_file" mapstructure:"config_file"`
	// FeedAddrs bscp feed server addresses
	FeedAddrs []string `json:"feed_addrs" mapstructure:"feed_addrs"`
	// FeedAddr bscp feed server address
	FeedAddr string `json:"feed_addr" mapstructure:"feed_addr"`
	// Biz bscp biz id
	Biz uint32 `json:"biz" mapstructure:"biz"`
	// Token bscp sdk token
	Token string `json:"token" mapstructure:"token"`
	// Apps bscp watched apps
	Apps []*AppConfig `json:"apps" mapstructure:"apps"`
	// Apps bscp watched app string
	App string `json:"app" mapstructure:"app"`
	// Labels bscp sdk labels
	Labels map[string]string `json:"labels" mapstructure:"labels"`
	// LabelsStr bscp sdk labels string
	LabelsStr string `json:"labels_str" mapstructure:"labels_str"`
	// UID bscp sdk uid
	UID string `json:"uid" mapstructure:"uid"`
	// TempDir config files temporary directory
	TempDir string `json:"temp_dir" mapstructure:"temp_dir"`
	// LabelsFile labels file path
	LabelsFile string `json:"labels_file" mapstructure:"labels_file"`
	// ConfigMatches global app config item's match conditions
	ConfigMatches []string `json:"config_matches" mapstructure:"config_matches"`
	// Port sidecar http server port
	Port int `json:"port" mapstructure:"port"`
	// EnableP2PDownload enable p2p download file
	EnableP2PDownload bool `json:"enable_p2p_download" mapstructure:"enable_p2p_download"`
	// BkAgentID bk gse agent id
	BkAgentID string `json:"bk_agent_id" mapstructure:"bk_agent_id"`
	// ClusterID bcs cluster id
	ClusterID string `json:"cluster_id" mapstructure:"cluster_id"`
	// PodID id of the pod where the bscp container resides
	PodID string `json:"pod_id" mapstructure:"pod_id"`
	// ContainerName bscp container name
	ContainerName string `json:"container_name" mapstructure:"container_name"`
	// FileCache file cache config
	FileCache *FileCacheConfig `json:"file_cache" mapstructure:"file_cache"`
	// KvCache kv cache config
	KvCache *KvCacheConfig `json:"kv_cache" mapstructure:"kv_cache"`
	// EnableMonitorResourceUsage 是否采集/监控资源使用率
	EnableMonitorResourceUsage bool `json:"enable_resource" mapstructure:"enable_resource"`
}

// String get config string
func (c *ClientConfig) String() string {
	conf := *c
	conf.Token = "******"
	cb, _ := json.Marshal(conf)
	return string(cb)
}

// ValidateBase validate the watch config
func (c *ClientConfig) ValidateBase() error {
	if len(c.FeedAddrs) == 0 {
		return fmt.Errorf("feed_addrs empty")
	}
	if c.Biz == 0 {
		return fmt.Errorf("biz is empty")
	}
	if c.Token == "" {
		return fmt.Errorf("token is empty")
	}

	if c.TempDir == "" {
		c.TempDir = constant.DefaultTempDir
	}
	if c.Port == 0 {
		c.Port = constant.DefaultHttpPort
	}
	if c.EnableP2PDownload {
		if c.BkAgentID == "" && (c.ClusterID == "" || c.PodID == "" || c.ContainerName == "") {
			return errors.New("to enable p2p download, either agent id must be set or cluster id, " +
				"pod id, container name must all be set")
		}
	}
	if c.FileCache == nil {
		c.FileCache = new(FileCacheConfig)
	}
	if err := c.FileCache.Validate(); err != nil {
		return err
	}
	if c.KvCache == nil {
		c.KvCache = new(KvCacheConfig)
	}
	if err := c.KvCache.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate validate the client config
func (c *ClientConfig) Validate() error {
	if err := c.ValidateBase(); err != nil {
		return err
	}
	if len(c.Apps) == 0 {
		return fmt.Errorf("apps is empty")
	}
	exists := make(map[string]bool)
	for _, app := range c.Apps {
		if exists[app.Name] {
			return fmt.Errorf("app %s is repeated ", app.Name)
		}
		if err := app.Validate(); err != nil {
			return err
		}
		exists[app.Name] = true
	}
	return nil
}

// AppConfig config for watched app
type AppConfig struct {
	// Name BSCP app name
	Name string `json:"name" mapstructure:"name"`
	// Labels instance labels
	Labels map[string]string `json:"labels" mapstructure:"labels"`
	// UID instance unique uid
	UID string `json:"uid" mapstructure:"uid"`
	// ConfigMatches app config item's match conditions
	ConfigMatches []string `json:"config_matches" mapstructure:"config_matches"`
}

// Validate validate the app watch config
func (c *AppConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("app is empty")
	}
	return nil
}

// FileCacheConfig config for file cache
type FileCacheConfig struct {
	// Enabled is whether enable file cache
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// CacheDir is file cache dir
	CacheDir string `json:"cache_dir" mapstructure:"cache_dir"`
	// cleanupIntervalSeconds is interval seconds of cleanup, not exposed for configuration now, use default value
	CleanupIntervalSeconds int64 `json:"-" mapstructure:"-"`
	// ThresholdGB is threshold gigabyte of cleanup
	ThresholdGB float64 `json:"threshold_gb" mapstructure:"threshold_gb"`
	// retentionRate is retention rate of cleanup, not exposed for configuration now, use default value
	RetentionRate float64 `json:"-" mapstructure:"-"`
}

// Validate validates the file cache config
func (c *FileCacheConfig) Validate() error {
	if c.CacheDir == "" {
		c.CacheDir = constant.DefaultFileCacheDir
	}
	if c.CleanupIntervalSeconds <= 0 {
		c.CleanupIntervalSeconds = constant.DefaultCleanupIntervalSeconds
	}
	if c.ThresholdGB <= 0 {
		c.ThresholdGB = constant.DefaultCacheThresholdGB
	}
	if c.RetentionRate <= 0 || c.RetentionRate > 1 {
		c.RetentionRate = constant.DefaultCacheRetentionRate
	}
	return nil
}

// KvCacheConfig config for kv cache
type KvCacheConfig struct {
	// Enabled is whether enable kv cache
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// ThresholdMB is threshold megabyte of kv cache
	ThresholdMB float64 `json:"threshold_mb" mapstructure:"threshold_mb"`
}

// Validate validates the kv cache config
func (c *KvCacheConfig) Validate() error {
	if c.ThresholdMB <= 0 {
		c.ThresholdMB = constant.DefaultKvCacheThresholdMB
	}
	return nil
}
