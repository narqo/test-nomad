job "be_server" {
    datacenters = ["dc-a", "dc-b"]

    # System type jobs run on every node, that satisfies the constraints.
    type = "system"

    # Restrict nodes to where the job can be deployed.
    constraint {
        attribute = "${meta.backend_role}"
        operator  = "set_contains"
        value     = "be_server"
    }

    update {
        auto_promote = false
        canary = 1
        max_parallel = 1
        health_check = "task_states"
        min_healthy_time = "5s"
    }

    meta = {
        version = "latest"
    }

    # ALL tasks within a group will be placed on the same host.
    group "server" {
        # Specify the number of these tasks we want (doesn't apply to job.type=system).
        count = 1

        # Task is an individual unit of work.
        task "server1" {
            driver = "raw_exec"
            config {
                command = "/home/vagrant/app/be_server/${NOMAD_META_VERSION}/bin/be_server"
                args = ["-http.addr", ":8081"]
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

        task "server2" {
            driver = "raw_exec"
            config {
                command = "/home/vagrant/app/be_server/${NOMAD_META_VERSION}/bin/be_server"
                args = ["-http.addr", ":8082"]
            }

            env {
                ADJUST_TAG = 2
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
