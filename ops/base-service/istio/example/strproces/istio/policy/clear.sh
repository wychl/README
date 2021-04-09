kubectl --insecure-skip-tls-verify delete virtualservice ui
kubectl --insecure-skip-tls-verify delete virtualservice processer
kubectl --insecure-skip-tls-verify delete gateway str-gateway
kubectl --insecure-skip-tls-verify delete destinationrule ui
kubectl --insecure-skip-tls-verify delete destinationrule processer
destinationrules virtualservices gateways
