package main

import (
	"encoding/json"

	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Application spec
type Application struct {
	Path     string       `json:",omitempty"`
	Gateway  *GatewaySpec `json:",omitempty"`
	Services []Service    `json:",omitempty"`
}

// GatewaySpec spec
type GatewaySpec struct {
	Name string `json:",omitempty"`
}

// Service spec
type Service struct {
	Name            string           `json:",omitempty"`
	Namespace       string           `json:",omitempty"`
	Gateways        []string         `json:",omitempty"`
	ServicePorts    []ServicePort    `json:",omitempty"`
	Deployments     []DeploymentSpec `json:",omitempty"`
	HTTPRoutes      []HTTPRoute      `json:",omitempty"`
	DestinationSpec *DestinationSpec `json:",omitempty"`
}

type DestinationSpec struct {
	Host          string         `json:",omitempty"`
	Subsets       []Subset       `json:",omitempty"`
	TrafficPolicy *TrafficPolicy `json:",omitempty"`
}

type Subset struct {
	Name          string            `json:",omitempty"`
	Labels        map[string]string `json:",omitempty"`
	TrafficPolicy *TrafficPolicy    `json:",omitempty"`
}

type TrafficPolicy struct {
	LoadBalancerType   string                             `json:",omitempty"`
	ConsistentHashType string                             `json:",omitempty"`
	LoadBalancerValue  interface{}                        `json:",omitempty"`
	ConnectionPool     *networking.ConnectionPoolSettings `json:",omitempty"`
}

type HTTPRoute struct {
	Matches               []Match
	HTTPRouteDestinations []HTTPRouteDestination
	HTTPRedirect          *networking.HTTPRedirect `json:",omitempty"`
	HTTPRewrite           *networking.HTTPRewrite  `json:",omitempty"`
	HTTPRetry             *networking.HTTPRetry    `json:",omitempty"`
	HTTPFaultInjection    *HTTPFaultInjection      `json:",omitempty"`
	Mirror                *Destination             `json:",omitempty"`
	CorsPolicy            *networking.CorsPolicy   `json:",omitempty"`
	Header                *Header                  `json:",omitempty"`
}

type Header struct {
	Request  *HeaderOperation `json:",omitempty"`
	Response *HeaderOperation `json:",omitempty"`
}

type HeaderOperation struct {
	Set    map[string]string `json:",omitempty"`
	Remove []string          `json:",omitempty"`
	Add    map[string]string `json:",omitempty"`
}

type HTTPFaultInjection struct {
	Delay *Delay `json:",omitempty"`
	Abort *Abort `json:",omitempty"`
}

type Delay struct {
	Percent   int32  `json:",omitempty"`
	DelayType string `json:",omitempty"`
	Duration  int64  `json:",omitempty"`
}

type Abort struct {
	Percent   int32       `json:",omitempty"`
	ErrorType string      `json:",omitempty"`
	Value     interface{} `json:",omitempty"`
}

type HTTPRouteDestination struct {
	Destination *Destination `json:",omitempty"`
	Weight      int32        `json:",omitempty"`
	Header      *Header      `json:",omitempty"`
}

type Match struct {
	Type         string            `json:",omitempty"`
	Way          string            `json:",omitempty"`
	Value        string            `json:",omitempty"`
	Headers      map[string]Match  `json:",omitempty"`
	SourceLabels map[string]string `json:",omitempty"`
	Gateways     []string          `json:",omitempty"`
	Port         uint32            `json:",omitempty"`
}

type Destination struct {
	Host   string
	Subset string
}

// ServicePort  service export port
type ServicePort struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
}

// DeploymentSpec spec
type DeploymentSpec struct {
	Name           string
	Version        string
	Replicas       int32
	Image          string
	ContainerPorts []int32
}

// ServiceSpec spec
type ServiceSpec struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   metav1.ObjectMeta `json:"metadata"`
	Spec       v1.ServiceSpec    `json:"spec"`
}

// VirtualService  spec
type VirtualService struct {
	APIVersion string                     `json:"apiVersion"`
	Kind       string                     `json:"kind"`
	Metadata   map[string]interface{}     `json:"metadata"`
	Spec       *networking.VirtualService `json:"spec"`
}

// MarshalJSON custom json marshal
func (s VirtualService) MarshalJSON() ([]byte, error) {
	jsonData, err := model.ToJSON(s.Spec)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Spec       json.RawMessage        `json:"spec"`
		APIVersion string                 `json:"apiVersion"`
		Kind       string                 `json:"kind"`
		Metadata   map[string]interface{} `json:"metadata"`
	}{
		Spec:       json.RawMessage(jsonData),
		APIVersion: s.APIVersion,
		Kind:       s.Kind,
		Metadata:   s.Metadata,
	})
}

// ServiceEntry spec
type ServiceEntry struct {
	APIVersion string                   `json:"apiVersion"`
	Kind       string                   `json:"kind"`
	Metadata   map[string]interface{}   `json:"metadata"`
	Spec       *networking.ServiceEntry `json:"spec"`
}

// MarshalJSON custom json marshal
func (s ServiceEntry) MarshalJSON() ([]byte, error) {
	jsonData, err := model.ToJSON(s.Spec)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Spec       json.RawMessage        `json:"spec"`
		APIVersion string                 `json:"apiVersion"`
		Kind       string                 `json:"kind"`
		Metadata   map[string]interface{} `json:"metadata"`
	}{
		Spec:       json.RawMessage(jsonData),
		APIVersion: s.APIVersion,
		Kind:       s.Kind,
		Metadata:   s.Metadata,
	})
}

// DestinationRule spec
type DestinationRule struct {
	APIVersion string                      `json:"apiVersion"`
	Kind       string                      `json:"kind"`
	Metadata   map[string]interface{}      `json:"metadata"`
	Spec       *networking.DestinationRule `json:"spec"`
}

// MarshalJSON custom json marshal
func (s DestinationRule) MarshalJSON() ([]byte, error) {
	jsonData, err := model.ToJSON(s.Spec)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Spec       json.RawMessage        `json:"spec"`
		APIVersion string                 `json:"apiVersion"`
		Kind       string                 `json:"kind"`
		Metadata   map[string]interface{} `json:"metadata"`
	}{
		Spec:       json.RawMessage(jsonData),
		APIVersion: s.APIVersion,
		Kind:       s.Kind,
		Metadata:   s.Metadata,
	})
}

// Gateway spec
type Gateway struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       *networking.Gateway    `json:"spec"`
}

// MarshalJSON custom json marshal
func (s Gateway) MarshalJSON() ([]byte, error) {
	jsonData, err := model.ToJSON(s.Spec)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Spec       json.RawMessage        `json:"spec"`
		APIVersion string                 `json:"apiVersion"`
		Kind       string                 `json:"kind"`
		Metadata   map[string]interface{} `json:"metadata"`
	}{
		Spec:       json.RawMessage(jsonData),
		APIVersion: s.APIVersion,
		Kind:       s.Kind,
		Metadata:   s.Metadata,
	})
}

// Deployment spec
type Deployment struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       appsv1.DeploymentSpec  `json:"spec"`
}
