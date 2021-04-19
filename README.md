# pigsty-cli


Command line tools for [pigsty](https://github.com/Vonng/pigsty)


## Man

```
NAME
    pigsty -- Pigsty Command-Line Interface v0.8 

SYNOPSIS               

    meta               setup meta nodes            init|fetch|repo|cache|ansible
    node               setup database nodes        init|ping|bash|ssh|admin 
    pgsql              setup postgres clusters     init|node|dcs|postgres|template|business|monitor|service|monly
    infra              setup infrastructure        init|ca|dns|prometheus|grafana|loki|haproxy|target
    clean              clean pgsql clusters        all|service|monitor|postgres|dcs
    config             mange pigsty config file    init|edit|info|dump|path
    serve              run pigsty API server       init|start|stop|restart|reload|status
    demo               setup local demo            init|up|new|clean|start|dns
    log                watch system log            query|postgres|patroni|pgbouncer|message
    pg                 pg operational tasks        user|db|svc|hba|log|psql|deploy|backup|restore|vacuum|repack


EXAMPLES

    1. infra summary
        pigsty infra

    2. pgsql clusters summary
        pigsty infra

    3. pigsty nodes summary
        pigsty node

    4. init pgsql cluster 'pg-test'
        pigsty pgsql init -l pg-test

    5. init new instance 10.10.10.13 of cluster 'pg-test'
        pigsty pgsql init -l 10.10.10.13

    6. remove cluster 'pg-test'
        pigsty clean -l pg-test

    7. create user dbuser_vonng on cluster 'pg-test'
        pigsty pg user dbuser_vonng -l pg-test

    8. create database test2 on cluster 'pg-test'
        pigsty pg db test -l pg-test
```


## Infra

```
infra -- setup pigsty infrastructure on meta node

    init           complete infra init on meta node
    repo           setup local yum repo
    ca             setup local ca
    dns            setup dnsmasq nameserver
    prometheus     setup prometheus & alertmanager
    grafana        setup grafana service
    loki           setup loki logging collector
    haproxy        refresh haproxy admin page index
    target         refresh prometheus static targets

Usage:
  pigsty infra [flags]
  pigsty infra [command]

Available Commands:
  ca          setup ca on meta node
  dns         setup dns infrastructure
  grafana     setup pigsty grafana on meta nodes
  haproxy     update haproxy index page
  init        init pigsty infra on meta nodes
  loki        setup loki on meta nodes
  node        setup node infrastructure
  pgsql       setup loki on meta nodes
  prometheus  setup pigsty prometheus on meta nodes
  repo        setup pigsty repo on meta nodes

Flags:
  -h, --help   help for infra

Global Flags:
  -i, --inventory string   inventory file (default "./pigsty.yml")
  -l, --limit string       limit execution hosts
  -t, --tags strings       limit execution tasks

Use "pigsty infra [command] --help" for more information about a command.
```




## Pgsql

```
SYNOPSIS:

    pgsql list                      show pgsql cluster definition
    pgsql init                      init new postgres clusters or instances
    pgsql node                      init pgsql node
    pgsql dcs                       init pgsql dcs (consul)
    pgsql postgres                  init postgres service (postgres|patroni|pgbouncer)
    pgsql monitor                   init monitor components
    pgsql service                   init services provider
    pgsql pgbouncer                 init pgbouncer service
    pgsql template                  init postgres template database
    pgsql business                  init postgres business users and databases
    pgsql config                    init patroni config template
    pgsql monly                     init monitor system in monitor-only mode
    pgsql hba                       init hba rule files
    pgsql remove                    remove postgres cluster or instances

Usage:
    pigsty pgsql [flags]
    pigsty pgsql [command]

Available Commands:
  business    init pgsql business user & db
  config      config pgsql with template
  dcs         init pgsql dcs service
  hba         init pgsql hba rules
  init        init pgsql on targets
  list        list pgsql clusters
  monitor     init pgsql monitor
  monly       init pgsql monitor only
  node        init pgsql node
  pgbouncer   init pgsql pgbouncer service
  postgres    init pgsql postgres service
  remove      remove pgsql from targets
  service     init pgsql service
  template    init pgsql template

Flags:
  -d, --detail   detail format
  -h, --help     help for pgsql
  -j, --json     json output
  -y, --yaml     yaml output

Global Flags:
  -i, --inventory string   inventory file (default "./pigsty.yml")
  -l, --limit string       limit execution hosts
  -t, --tags strings       limit execution tasks

Use "pigsty pgsql [command] --help" for more information about a command.
```