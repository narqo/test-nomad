job "be_worker" {
    datacenters = ["dc-a", "dc-b"]

    # System type jobs run on every node, that satisfies the constraints.
    type = "system"

    # Restrict nodes to where the job can be deployed.
    constraint {
        attribute = "${meta.backend_role}"
        operator  = "set_contains"
        value     = "be_worker"
    }

    update {
        max_parallel = 1
        stagger = "10s"
    }

    meta = {
        version = "latest"
    }

    # ALL tasks within a group will be placed on the same host.
    group "dc_a_worker" {
        constraint {
            attribute = "${node.datacenter}"
            value     = "dc-a"
        }

        # Task is an individual unit of work.
        task "worker1" {
            driver = "raw_exec"
            config {
                command = "/home/vagrant/app/be_worker/${NOMAD_META_VERSION}/bin/be_worker"
            }

            env {
                ADJUST_TAG = 1
                ADJUST_NODE_ID = "${node.unique.id}"
                ADJUST_NODE_NAME = "${node.unique.name}"
                GO_ENV = "${NOMAD_DC}_production"
                GO_CONFIG = "config.yml"
            }

            resources {
                memory = 128
            }
        }

        task "worker2" {
            driver = "raw_exec"
            config {
                command = "/home/vagrant/app/be_worker/${NOMAD_META_VERSION}/bin/be_worker"
            }

            env {
                ADJUST_TAG = 1
                ADJUST_NODE_ID = "${node.unique.id}"
                ADJUST_NODE_NAME = "${node.unique.name}"
                GO_ENV = "${NOMAD_DC}_production"
                GO_CONFIG = "config.yml"
            }

            resources {
                memory = 128
            }
        }
    }

    group "dc_b_worker" {
        constraint {
            attribute = "${node.datacenter}"
            value     = "dc-b"
        }

        task "worker1" {
            driver = "raw_exec"
            config {
                command = "/home/vagrant/app/be_worker/${NOMAD_META_VERSION}/bin/be_worker"
            }

            env {
                ADJUST_TAG = 1
                ADJUST_NODE_ID = "${node.unique.id}"
                ADJUST_NODE_NAME = "${node.unique.name}"
                GO_ENV = "${NOMAD_DC}_production"
                GO_CONFIG = "config.yml"
            }

            resources {
                memory = 128
            }
        }
    }
}
