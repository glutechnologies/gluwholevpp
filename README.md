# Gluwholevpp

Software written in Go used for generating Vlan transformations using VPP dataplane in order to provide wholesale bitstream services

## Forward API Unix Socket
In order to develop this control plane sometimes is useful to forward VPP Unix socket from vpp device to a development machine. We can use SSH forwarding capabilities:

```
ssh root@<vpp-management-ip> -L<local-sock>:/run/vpp/cli.sock
```