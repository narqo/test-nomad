# RFC: Private Cloud with Nomad (Test)

## Inventory

DC A + DC B (5 nodes)

- dc-a-node-[1,3]
- dc-b-node-[1,2]

### Nomad

Server `dc-a-node-1`; clients `dc-a-node-[1,3]` + `dc-b-node-[1,2]`.

Every client node is annotated with roles, that define the type of workload the node allowed to run. I.e.
`backend_role=be_server,be_worker`.

A service (nomad job) uses the role as a constraints:

```
// be_server.nomad

job "be_server {
  // the job is allowed to run only on the nodes annotated with "backend_role=be_server"
  constraint {
    attribute = "${meta.backend_role}"
    operator  = "set_contains"
    value     = "be_server"
  }
}
```

#### Haproxy

LB `dc-b-node-2`

### Provisioning

```
$ vagrant provision [node]
```

See `Vagrantfile` and `site.yml`.

## Stage 0

### Workload

#### DC A

Nodes `dc-a-node-[1,2]` run

- `be_server` - 2 instances per node
- local Redis queue

Nodes `dc-a-node-3` run

- `be_worker` - 1 instance per corresponding `be_server` in the DC (2 instances)

#### DC B

Nodes `dc-b-node-1` run

- `be_server` - 2 instances per node
- local Redis queue

- `be_worker` - 1 instance per corresponding `be_server` in the DC

Nodes `dc-b-node-2` run

- `be_worker` - 1 instance per corresponding `be_server` in the DC (1 instance)

### Deployment

```
$ ansible-playbook --inventory inventory/10-hosts \
  [--private-key ~/.vagrant.d/insecure_private_key] [--tags be_server] [--limit <node>] deploy.yml
```

# Extra

Test hosts provisioning:

```
$ ansible --inventory inventory/10-hosts --private-key ~/.vagrant.d/insecure_private_key -m ping all
```

Test be_server, through the LB

```
$ while true ; curl -s 'http://172.16.2.102:8080/internal/stats' ; sleep 0.5; end
```
