local gateway = import 'loki/gateway.libsonnet';
local loki = import 'loki/loki.libsonnet';
local promtail = import 'promtail/promtail.libsonnet';

loki + promtail + gateway {

  distributor_deployment+: {
      replicas: 1,
  },
  gateway_deployment+:{
      replicas: 1,
  },
  ingester_deployment+:{
      replicas: 1,
  },
  querier_deployment+:{
      replicas: 1,
  },
  table_manager_deployment:{},
  table_manager_service:{},
  consul_deployment:{},
  _config+:: {
    namespace: 'loki',
    htpasswd_contents: 'loki:$apr1$H4yGiGNg$ssl5/NymaGFRUvxIV1Nyr.',
    bigtable_instance: '',
    bigtable_project:  '',
    gcs_bucket_name:  '',
    storage_backend: 'boltdb',
    replication_factor: 1,


    promtail_config+: {
      scheme: 'http',
      hostname: 'gateway.%(namespace)s.svc' % $._config,
      username: 'loki',
      password: 'password',
      container_root_path: '/var/lib/docker',
    },

    loki+: {
        storage_config: {
            boltdb: {
                directory: '/tmp/loki/index',
            },
            filesystem: {
                directory: '/tmp/loki/chunks',
            }
        },
        ingester: {
            chunk_idle_period: '15m',
            chunk_block_size: 262144,

            lifecycler: {
                address: '127.0.0.1',
                ring: {
                    store: 'inmemory',
                    replication_factor: 1,
                },
            },
        },
        schema_config: {
            configs: [{
                    from: '0',
                    store: 'boltdb',
                    object_store: 'filesystem',
                    schema: 'v9',
                    index: {
                        prefix: '%s_index_' % $._config.table_prefix,
                        period: '168h',
                    },
            }],
        },
    },
  },
}