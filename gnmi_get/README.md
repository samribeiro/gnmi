# gNMI Get

A simple shell binary that performs a GET against a gNMI Target.

## Install

```
go get github.com/samribeiro/gnmi/gnmi_get
go install github.com/samribeiro/gnmi/gnmi_get
```

## Run

```
gnmi_get \
  -target_address localhost:32123 \
  -key client.key \
  -cert client.crt \
  -ca ca.crt \
  -target_name server \
  -query "system/openflow/controllers/controller[main]/connections/connection[0]/state/address"
```
