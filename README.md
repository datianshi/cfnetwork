# cfnetwork

## curlv2

This tries to identify network connectivity for cloudfoundry deployment

### Curl v2 through load balancer

```
cfnetwork curlv2 --domain 10.244.0.34.xip.io
```

```
cfnetwork curlv2 --domain 10.244.0.34.xip.io --https
```

### Curl v2 through router directly

```
cfnetwork curlv2 router --domain 10.244.0.34.xip.io --ip 10.244.0.22
```
