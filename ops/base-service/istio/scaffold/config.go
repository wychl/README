package main

import (
	"fmt"

	google_protobuf "github.com/gogo/protobuf/types"
	networking "istio.io/api/networking/v1alpha3"
)

// Config  config
type Config struct {
	Path     string
	Gateway  string
	Services []ServiceConfig
}

type ServiceConfig struct {
	Name         string
	Gateway      string
	ServicePorts []ServicePort
	Containers   []ContainerConfig
}

// Container config
type ContainerConfig struct {
	Name    string
	Version string
	Image   string
	Ports   []int32
}

func (c *Config) ToApplication() *Application {
	app := Application{
		Path: c.Path,
	}

	var services []Service
	for _, item := range c.Services {
		var s Service
		s.Name = item.Name
		s.Gateways = []string{item.Gateway}
		s.ServicePorts = item.ServicePorts

		s.HTTPRoutes = append(s.HTTPRoutes)

		for _, container := range item.Containers {
			s.HTTPRoutes = append(s.HTTPRoutes, container.httpRoute())
			s.Deployments = append(s.Deployments, container.deployment())
		}

		s.DestinationSpec = item.toDestination()
		services = append(services, s)
	}
	app.Services = services

	if c.Gateway != "" {
		app.Gateway = &GatewaySpec{
			Name: c.Gateway,
		}
	}

	return &app
}

func (c *ContainerConfig) httpRoute() HTTPRoute {
	r := HTTPRoute{
		Matches: []Match{
			Match{
				Type:  "uri",
				Way:   "exact",
				Value: c.Version,
			},
		},
		HTTPRewrite: &networking.HTTPRewrite{
			Uri: "/",
		},
		HTTPRouteDestinations: []HTTPRouteDestination{
			HTTPRouteDestination{
				Destination: &Destination{
					Host:   c.Name,
					Subset: "v1",
				},
				Weight: 100,
			},
		},
	}

	return r
}

func (c *ContainerConfig) deployment() DeploymentSpec {
	r := DeploymentSpec{
		Name:           fmt.Sprintf("%s-%s", c.Name, c.Version),
		Version:        c.Version,
		ContainerPorts: c.Ports,
		Replicas:       1,
		Image:          c.Image,
	}

	return r
}

func (s *ServiceConfig) toDestination() *DestinationSpec {
	r := &DestinationSpec{
		Host: s.Name,
	}

	var subsets []Subset
	for _, item := range s.Containers {
		var subset Subset
		subset.Name = item.Version
		subset.Labels = map[string]string{"version": item.Version}
		subset.TrafficPolicy = &TrafficPolicy{
			LoadBalancerType:  "simple",
			LoadBalancerValue: 0,
			ConnectionPool: &networking.ConnectionPoolSettings{
				Tcp: &networking.ConnectionPoolSettings_TCPSettings{
					MaxConnections: 5,
					ConnectTimeout: &google_protobuf.Duration{Seconds: 30},
					TcpKeepalive: &networking.ConnectionPoolSettings_TCPSettings_TcpKeepalive{
						Time: &google_protobuf.Duration{
							Seconds: 5,
						},
					},
				},
			},
		}
		subsets = append(subsets, subset)
	}
	r.Subsets = subsets

	return r
}
