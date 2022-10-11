package model

import (
	v1 "k8s.io/api/apps/v1"
)

type Deployments struct {
	Items []*v1.Deployment `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ReqReplica struct {
	NameSpace string `json:"ns"`
	Deploy    string `json:"deploy"`
	Dec       bool   `json:"dec"`
}

type DeployDetail struct {
	Name       string   `json:"name"`
	ImageNames []string `json:"images_name"`
	Replicas   []int32  `json:"replicas"` // index表示表示不同的含义，0 总，1 可用，2 不可用
	Pods       []*Pod   `json:"pods"`
}

type Pod struct {
	Name         string   `json:"name"`
	Ns           string   `json:"ns"`
	Images       []string `json:"images"`
	NodeName     string   `json:"node_name"`
	CreateTime   string   `json:"create_time"`
	RestartCount int32    `json:"restart_count"`
	Phase        string   `json:"phase"`
	Message      string   `json:"message"`
	EventMsg     string   `json:"event_msg"`
	Ready        bool     `json:"ready"`
}
