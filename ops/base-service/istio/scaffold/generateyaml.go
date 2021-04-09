package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/ghodss/yaml"
	google_protobuf "github.com/gogo/protobuf/types"
	networking "istio.io/api/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GenerateYaml generate istio deploy file
func (a *Application) GenerateYaml() error {
	f, err := os.OpenFile(a.Path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer

	gatewayYaml, err := a.Gateway.generateGatewayYaml()
	if err != nil {
		return err
	}
	_, err = buf.WriteString(gatewayYaml)
	if err != nil {
		return err
	}
	_, err = buf.WriteString("\n---\n")
	if err != nil {
		return err
	}

	for _, service := range a.Services {
		virtualServiceYaml, err := service.generateVirtualServiceYaml()
		if err != nil {
			return err
		}
		_, err = buf.WriteString(virtualServiceYaml)
		if err != nil {
			return err
		}
		_, err = buf.WriteString("\n---\n")
		if err != nil {
			return err
		}

		destinationRuleYaml, err := service.generateDestinationRule()
		if err != nil {
			return err
		}
		_, err = buf.WriteString(destinationRuleYaml)
		if err != nil {
			return err
		}
		_, err = buf.WriteString("\n---\n")
		if err != nil {
			return err
		}

		serviceYaml, err := service.generateServiceYaml()
		if err != nil {
			return err
		}
		_, err = buf.WriteString(serviceYaml)
		if err != nil {
			return err
		}
		_, err = buf.WriteString("\n---\n")
		if err != nil {
			return err
		}

		deployYaml, err := service.generateDeployYaml()
		if err != nil {
			return err
		}
		_, err = buf.WriteString(deployYaml)
		if err != nil {
			return err
		}
		_, err = buf.WriteString("\n---\n")
		if err != nil {
			return err
		}
	}

	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) generateVirtualServiceYaml() (string, error) {
	virtualService := VirtualService{
		APIVersion: "networking.istio.io/v1alpha3",
		Kind:       "VirtualService",
		Metadata: map[string]interface{}{
			"name": s.Name,
		},
		Spec: &networking.VirtualService{
			Hosts:    []string{"*"},
			Gateways: s.Gateways,
		},
	}

	var routes []*networking.HTTPRoute
	for _, route := range s.HTTPRoutes {
		var matches []*networking.HTTPMatchRequest
		for _, match := range route.Matches {
			hmr := &networking.HTTPMatchRequest{}
			switch match.Type {
			case uri:
				hmr.Uri = convertTOStringMatch(match.Way, match.Value)
			case scheme:
				hmr.Scheme = convertTOStringMatch(match.Way, match.Value)
			case method:
				hmr.Method = convertTOStringMatch(match.Way, match.Value)
			case authority:
				hmr.Authority = convertTOStringMatch(match.Way, match.Value)
			case headers:
				hmr.Headers = make(map[string]*networking.StringMatch)
				for key, value := range match.Headers {
					hmr.Headers[key] = convertTOStringMatch(value.Way, value.Value)
				}
			case port:
				hmr.Port = match.Port
			case sourceLabels:
				hmr.SourceLabels = match.SourceLabels
			case gateways:
				hmr.Gateways = match.Gateways

			}
			matches = append(matches, hmr)
		}

		var destinations []*networking.HTTPRouteDestination
		for _, hrd := range route.HTTPRouteDestinations {
			var nhrd networking.HTTPRouteDestination
			nhrd.Weight = hrd.Weight
			nhrd.Destination = &networking.Destination{
				Host:   hrd.Destination.Host,
				Subset: hrd.Destination.Subset,
			}

			if hrd.Header != nil {
				// nhrd.Headers = &networking.Headers{}
				// if hrd.Header.Request != nil {
				// 	nhrd.Headers.Request = &networking.Headers_HeaderOperations{
				// 		Set:    hrd.Header.Request.Set,
				// 		Add:    hrd.Header.Request.Add,
				// 		Remove: hrd.Header.Request.Remove,
				// 	}
				// }
				// if hrd.Header.Response != nil {
				// 	nhrd.Headers.Request = &networking.Headers_HeaderOperations{
				// 		Set:    hrd.Header.Response.Set,
				// 		Add:    hrd.Header.Response.Add,
				// 		Remove: hrd.Header.Response.Remove,
				// 	}
				// }
			}
			destinations = append(destinations, &nhrd)
		}

		httpRoute := &networking.HTTPRoute{
			Match:      matches,
			Route:      destinations,
			Redirect:   route.HTTPRedirect,
			Rewrite:    route.HTTPRewrite,
			Retries:    route.HTTPRetry,
			CorsPolicy: route.CorsPolicy,
		}

		if route.Header != nil {
			httpRoute.Headers = &networking.Headers{}
			if route.Header.Request != nil {
				httpRoute.Headers.Request = &networking.Headers_HeaderOperations{
					Set:    route.Header.Request.Set,
					Add:    route.Header.Request.Add,
					Remove: route.Header.Request.Remove,
				}
			}
			if route.Header.Response != nil {
				httpRoute.Headers.Request = &networking.Headers_HeaderOperations{
					Set:    route.Header.Response.Set,
					Add:    route.Header.Response.Add,
					Remove: route.Header.Response.Remove,
				}
			}
		}

		if route.Mirror != nil {
			httpRoute.Mirror = &networking.Destination{
				Host:   route.Mirror.Host,
				Subset: route.Mirror.Subset,
			}
		}

		if route.HTTPFaultInjection != nil {
			var fault networking.HTTPFaultInjection
			if route.HTTPFaultInjection.Delay != nil {
				switch route.HTTPFaultInjection.Delay.DelayType {
				case fixedDelay:
					fault.Delay = &networking.HTTPFaultInjection_Delay{
						Percent: route.HTTPFaultInjection.Delay.Percent,
						HttpDelayType: &networking.HTTPFaultInjection_Delay_FixedDelay{
							FixedDelay: &google_protobuf.Duration{
								Seconds: route.HTTPFaultInjection.Delay.Duration,
							},
						},
					}
				case exponentialDelay:
					fault.Delay = &networking.HTTPFaultInjection_Delay{
						Percent: route.HTTPFaultInjection.Delay.Percent,
						HttpDelayType: &networking.HTTPFaultInjection_Delay_ExponentialDelay{
							ExponentialDelay: &google_protobuf.Duration{
								Seconds: route.HTTPFaultInjection.Delay.Duration,
							},
						},
					}
				}

			}
			if route.HTTPFaultInjection.Abort != nil {
				switch route.HTTPFaultInjection.Abort.ErrorType {
				case abortHttp:
					fault.Abort = &networking.HTTPFaultInjection_Abort{
						Percent: route.HTTPFaultInjection.Abort.Percent,
						ErrorType: &networking.HTTPFaultInjection_Abort_HttpStatus{
							HttpStatus: int32(route.HTTPFaultInjection.Abort.Value.(int)),
						},
					}
				case abortGrpc:
					fault.Abort = &networking.HTTPFaultInjection_Abort{
						Percent: route.HTTPFaultInjection.Abort.Percent,
						ErrorType: &networking.HTTPFaultInjection_Abort_GrpcStatus{
							GrpcStatus: route.HTTPFaultInjection.Abort.Value.(string),
						},
					}
				case abortHttp2Error:
					fault.Abort = &networking.HTTPFaultInjection_Abort{
						Percent: route.HTTPFaultInjection.Abort.Percent,
						ErrorType: &networking.HTTPFaultInjection_Abort_Http2Error{
							Http2Error: route.HTTPFaultInjection.Abort.Value.(string),
						},
					}
				}
			}
			httpRoute.Fault = &fault

		}

		routes = append(routes, httpRoute)

	}

	virtualService.Spec.Http = routes
	d, err := yaml.Marshal(virtualService)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func (g GatewaySpec) generateGatewayYaml() (string, error) {
	gateway := Gateway{
		APIVersion: "networking.istio.io/v1alpha3",
		Kind:       "Gateway",
		Metadata: map[string]interface{}{
			"name": g.Name,
		},
		Spec: &networking.Gateway{
			Selector: map[string]string{"istio": "ingressgateway"},
			Servers: []*networking.Server{&networking.Server{
				Port: &networking.Port{
					Number:   80,
					Name:     "http",
					Protocol: "http",
				},
				Hosts: []string{"*"},
			}},
		},
	}

	d, err := yaml.Marshal(gateway)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

func (s Service) generateServiceYaml() (string, error) {
	service := ServiceSpec{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: metav1.ObjectMeta{
			Name: s.Name,
			Labels: map[string]string{
				"app":     s.Name,
				"service": s.Name,
			},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": s.Name},
		},
	}

	var servicePorts []v1.ServicePort
	for _, port := range s.ServicePorts {
		servicePorts = append(servicePorts, v1.ServicePort{
			Name:     port.Name,
			Protocol: v1.Protocol(port.Protocol),
			Port:     port.Port,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: port.TargetPort,
			},
		})
	}
	service.Spec.Ports = servicePorts

	d, err := yaml.Marshal(service)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

func (s Service) generateDeployYaml() (string, error) {
	var result string

	for _, value := range s.Deployments {
		var containerPorts []v1.ContainerPort
		for _, port := range value.ContainerPorts {
			containerPorts = append(containerPorts, v1.ContainerPort{
				ContainerPort: port,
			})
		}

		deploy := Deployment{
			APIVersion: "extensions/v1beta1",
			Kind:       "Deployment",
			Metadata:   map[string]interface{}{"name": value.Name},
			Spec: appsv1.DeploymentSpec{
				Replicas: &value.Replicas,
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":     s.Name,
							"version": value.Version,
						},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{v1.Container{
							Name:            s.Name,
							Image:           value.Image,
							ImagePullPolicy: v1.PullAlways,
							Ports:           containerPorts,
						}},
					},
				},
			},
		}

		d, err := yaml.Marshal(deploy)
		if err != nil {
			return "", err
		}

		if result != "" {
			result = fmt.Sprintf("%s\n---\n%s", result, d)
		} else {
			result = string(d)
		}

	}

	return result, nil
}

func (s Service) generateDestinationRule() (string, error) {

	destinationRule := DestinationRule{
		APIVersion: "networking.istio.io/v1alpha3",
		Kind:       "DestinationRule",
		Metadata: map[string]interface{}{
			"name": s.Name,
		},
		Spec: &networking.DestinationRule{
			Host: s.Name,
		},
	}

	var subsets []*networking.Subset
	for _, subset := range s.DestinationSpec.Subsets {
		s := &networking.Subset{
			Name:   subset.Name,
			Labels: subset.Labels,
			TrafficPolicy: &networking.TrafficPolicy{
				LoadBalancer: &networking.LoadBalancerSettings{},
			},
		}

		switch subset.TrafficPolicy.LoadBalancerType {
		case "simple":
			s.TrafficPolicy.LoadBalancer = &networking.LoadBalancerSettings{
				LbPolicy: &networking.LoadBalancerSettings_Simple{
					Simple: networking.LoadBalancerSettings_SimpleLB(subset.TrafficPolicy.LoadBalancerValue.(int)),
				},
			}
		case "consistentHash":
			switch subset.TrafficPolicy.ConsistentHashType {
			case "httpCookie":
				value := subset.TrafficPolicy.LoadBalancerValue.(map[string]interface{})
				duration := time.Duration(value["Ttl"].(int64))

				s.TrafficPolicy.LoadBalancer = &networking.LoadBalancerSettings{
					LbPolicy: &networking.LoadBalancerSettings_ConsistentHash{
						ConsistentHash: &networking.LoadBalancerSettings_ConsistentHashLB{
							HashKey: &networking.LoadBalancerSettings_ConsistentHashLB_HttpCookie{
								HttpCookie: &networking.LoadBalancerSettings_ConsistentHashLB_HTTPCookie{
									Name: value["name"].(string),
									Path: value["Path"].(string),
									Ttl:  &duration,
								},
							},
						},
					},
				}
			case "httpHeader":
				s.TrafficPolicy.LoadBalancer = &networking.LoadBalancerSettings{
					LbPolicy: &networking.LoadBalancerSettings_ConsistentHash{
						ConsistentHash: &networking.LoadBalancerSettings_ConsistentHashLB{
							HashKey: &networking.LoadBalancerSettings_ConsistentHashLB_HttpHeaderName{
								HttpHeaderName: subset.TrafficPolicy.LoadBalancerValue.(string),
							},
						},
					},
				}
			case "useSourceIp":
				s.TrafficPolicy.LoadBalancer = &networking.LoadBalancerSettings{
					LbPolicy: &networking.LoadBalancerSettings_ConsistentHash{
						ConsistentHash: &networking.LoadBalancerSettings_ConsistentHashLB{
							HashKey: &networking.LoadBalancerSettings_ConsistentHashLB_UseSourceIp{
								UseSourceIp: subset.TrafficPolicy.LoadBalancerValue.(bool),
							},
						},
					},
				}
			}

		}
		subsets = append(subsets, s)
	}

	destinationRule.Spec.Subsets = subsets

	if s.DestinationSpec.TrafficPolicy != nil {
		switch s.DestinationSpec.TrafficPolicy.LoadBalancerType {
		case "simple":
			destinationRule.Spec.TrafficPolicy.LoadBalancer = &networking.LoadBalancerSettings{
				LbPolicy: &networking.LoadBalancerSettings_Simple{
					Simple: networking.LoadBalancerSettings_SimpleLB(s.DestinationSpec.TrafficPolicy.LoadBalancerValue.(int)),
				},
			}
		case "consistentHash":
			switch s.DestinationSpec.TrafficPolicy.ConsistentHashType {
			case "httpCookie":
				value := s.DestinationSpec.TrafficPolicy.LoadBalancerValue.(map[string]interface{})
				duration := time.Duration(value["Ttl"].(int64))
				destinationRule.Spec.TrafficPolicy.LoadBalancer = &networking.LoadBalancerSettings{
					LbPolicy: &networking.LoadBalancerSettings_ConsistentHash{
						ConsistentHash: &networking.LoadBalancerSettings_ConsistentHashLB{
							HashKey: &networking.LoadBalancerSettings_ConsistentHashLB_HttpCookie{
								HttpCookie: &networking.LoadBalancerSettings_ConsistentHashLB_HTTPCookie{
									Name: value["name"].(string),
									Path: value["Path"].(string),
									Ttl:  &duration,
								},
							},
						},
					},
				}
			case "httpHeader":
				destinationRule.Spec.TrafficPolicy.LoadBalancer = &networking.LoadBalancerSettings{
					LbPolicy: &networking.LoadBalancerSettings_ConsistentHash{
						ConsistentHash: &networking.LoadBalancerSettings_ConsistentHashLB{
							HashKey: &networking.LoadBalancerSettings_ConsistentHashLB_HttpHeaderName{
								HttpHeaderName: s.DestinationSpec.TrafficPolicy.LoadBalancerValue.(string),
							},
						},
					},
				}
			case "useSourceIp":
				destinationRule.Spec.TrafficPolicy.LoadBalancer = &networking.LoadBalancerSettings{
					LbPolicy: &networking.LoadBalancerSettings_ConsistentHash{
						ConsistentHash: &networking.LoadBalancerSettings_ConsistentHashLB{
							HashKey: &networking.LoadBalancerSettings_ConsistentHashLB_UseSourceIp{
								UseSourceIp: s.DestinationSpec.TrafficPolicy.LoadBalancerValue.(bool),
							},
						},
					},
				}
			}

		}
	}

	d, err := yaml.Marshal(destinationRule)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
