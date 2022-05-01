# Gluwholevpp

Software written in Go used for generating Vlan transformations using VPP dataplane in order to provide wholesale bitstream services

## Forward API Unix Socket
In order to develop this control plane sometimes is useful to forward VPP Unix socket from vpp device to a development machine. We can use SSH forwarding capabilities:

```
ssh root@<vpp-management-ip> -L<local-sock>:/run/vpp/cli.sock
```

## Installation

```
# Build gluwholevpp
go build -o bin

# Copy systemd unit from misc/gluwholevpp.service
cp misc/gluwholevpp.service /etc/systemd/system/

# Copy software to /opt/glutec/gluwholevpp
mkdir -p /opt/glutec/gluwholevpp
cp bin/gluwholevpp /opt/glutec/gluwholevpp
cp data.db /opt/glutec/gluwholevpp
cp gluwholevpp.default.toml /opt/glutec/gluwholevpp/gluwholevpp.toml

systemctl daemon-reload
systemctl enable gluwholevpp
systemctl start gluwholevpp
```