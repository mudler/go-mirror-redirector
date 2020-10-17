# [![Docker Repository on Quay](https://quay.io/repository/mocaccino/mirror-redirector/status "Docker Repository on Quay")](https://quay.io/repository/mocaccino/mirror-redirector) go-mirror-redirector

Simple mirror redirector in golang.

No frills geo-ip based mirror redirector, supports a yaml file as config, in the following format:

```yaml
COUNTRYCODE:
- url1
- url2
AU:
- url2
IT:
- foo

default:
- defaultfallback1
- defaultfallback2
```

you can pass by a config path with `CONFIG` environment variable. You can customize a listening address with `HOST`, the port with `PORT`, and the deployment mode (prod, dev) with `MACARON_ENV`.

There is also a deployment example for Kubernetes in [kube.yaml](https://github.com/mudler/go-mirror-redirector/blob/main/kube.yaml).