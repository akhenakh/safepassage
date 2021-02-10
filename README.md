# safepassage
A Drone plugin to safely extract secrets.

## Disclaimer
This is not a nice solution for your lack of discipline, your secrets should be stored elsewhere.  
When recovered please rotate your secret.

## How to use safepassage
safepassage is a simple plugins that will export your secret into an OpenPGP encrypted message.

Provide your pubkey, create a temporary branch change the `.drone.yml` using the following example and commit to trigger a build.  
Secrets can be passed by environment or by settings.

You can also specify a "format" setting. Accepted values are `std` or `env` (which base64 encodes the values and presents them in an env-file format). The default is `std`.

```yaml
kind: pipeline

steps:
- name: extractor
  image: akhenakh/safepassage:1.1
  environment:
    my_secret:
      from_secret: my_drone_secret
  settings:
    another_secret:
      from_secret: my_other_drone_secret
    secrets:
      - MY_SECRET
      - ANOTHER_SECRET
    pubkey: |
      -----BEGIN PGP PUBLIC KEY BLOCK-----
      ....
      -----END PGP PUBLIC KEY BLOCK-----
```

## Details
safepassage is a simple Go binary build into a Distroless Docker image, it used [GopenPGP](https://gopenpgp.org/) implementation.

## Background
```shell script
MY_SECRET=hello && ./safepassage -secrets=MY_SECRET -pubKey="$(cat testdata/pubkey.asc)"      
```

Env are prefixed by `PLUGIN_`.
