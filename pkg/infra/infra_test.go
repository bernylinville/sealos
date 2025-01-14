// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package infra

import (
	"fmt"
	"testing"
	"time"

	"github.com/fanux/sealos/pkg/infra/aliyun"
	"github.com/fanux/sealos/pkg/infra/huawei"

	"sigs.k8s.io/yaml"

	v2 "github.com/fanux/sealos/pkg/types/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAliApply(t *testing.T) {
	//setup infra
	infra := v2.Infra{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Infra",
			APIVersion: v2.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-infra",
		},
		Spec: v2.InfraSpec{
			Cluster: v2.Cluster{
				RegionIDs: []string{"cn-shanghai"},
				ZoneIDs:   []string{"cn-shanghai-l"},
				AccessChannels: v2.AccessChannels{
					SSH: v2.SSH{
						Passwd: "Fanux#123",
						Port:   22,
					},
				},
				Metadata: v2.ClusterMeta{
					IsSeize: true,
				},
			},
			Hosts: []v2.Host{
				{
					Roles:  []string{"master", "ssdxxx"},
					CPU:    2,
					Memory: 4,
					Count:  1,
					Disks:  []v2.Disk{},
					OS: v2.OS{
						Name: "ubuntu",
					},
				},
			},
			Provider: aliyun.AliyunProvider,
		},
	}

	aliProvider, err := NewDefaultProvider(&infra)
	if err != nil {
		fmt.Printf("%v", err)
	} else {
		fmt.Printf("%v", aliProvider.Apply())
	}
	// todo
	t.Run("modify instance system disk", func(t *testing.T) {
		j, _ := yaml.Marshal(&infra)
		t.Log("output yaml:", string(j))
		infra.Spec.Hosts = []v2.Host{
			{
				Roles:  []string{"master", "ssd"},
				CPU:    2,
				Memory: 4,
				Count:  1,
				Disks:  []v2.Disk{},
				OS: v2.OS{
					Name: "centos",
				},
			},
			{
				Roles:  []string{"master", "ssdxxx"},
				CPU:    2,
				Memory: 4,
				Count:  1,
				Disks:  []v2.Disk{},
				OS: v2.OS{
					Name: "ubuntu",
				},
			},
		}
		t.Log(fmt.Sprintf("add server:%v", aliProvider.Apply()))
		j, _ = yaml.Marshal(&infra)
		t.Log("output yaml:", string(j))
		time.Sleep(10 * time.Second)
		infra.Spec.Hosts = []v2.Host{
			{
				Roles:  []string{"master", "ssd"},
				CPU:    2,
				Memory: 4,
				Count:  1,
				Disks:  []v2.Disk{},
				OS: v2.OS{
					Name: "centos",
				},
			},
		}
		t.Log(fmt.Sprintf("delete:%v", aliProvider.Apply()))
		j, _ = yaml.Marshal(&infra)
		t.Log("output yaml:", string(j))
	})
	//teardown
	time.Sleep(20 * time.Second)
	now := metav1.Now()
	infra.ObjectMeta.DeletionTimestamp = &now
	t.Log(fmt.Sprintf("%v", aliProvider.Apply()))
}

func TestHuaweiApply(t *testing.T) {
	//setup infra
	infra := v2.Infra{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Infra",
			APIVersion: v2.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-infra",
		},
		Spec: v2.InfraSpec{
			Credential: v2.Credential{ProjectID: "06b264130780105b2f50c0145ba32d41"},
			Cluster: v2.Cluster{
				RegionIDs: []string{"cn-north-4"},
				ZoneIDs:   []string{""},
				AccessChannels: v2.AccessChannels{
					SSH: v2.SSH{
						Passwd: "Fanux#123",
						Port:   22,
					},
				},
				Metadata: v2.ClusterMeta{IsSeize: true},
			},
			Hosts: []v2.Host{
				{
					Roles:  []string{"master", "ssdxxx"},
					CPU:    2,
					Memory: 4,
					Count:  1,
					Disks:  []v2.Disk{},
					OS: v2.OS{
						Name: "ubuntu",
					},
				},
			},
			Provider: huawei.HuaweiProvider,
		},
	}

	hwProvider, err := NewDefaultProvider(&infra)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", hwProvider.Apply())
	// todo
	t.Run("default", func(t *testing.T) {

	})
	//teardown
	time.Sleep(20 * time.Second)
	now := metav1.Now()
	infra.ObjectMeta.DeletionTimestamp = &now
	t.Log(fmt.Sprintf("%v", hwProvider.Apply()))
}
