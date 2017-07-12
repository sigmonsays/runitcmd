# runitcmd

manage runit

# features
- manage multiple service matching via patterns
- easily create and manage runit service
- import/export service configuration
- user and global configuration for setting common parameters

# install

Assuming you have a working go environment

    go get -u github.com/sigmonsays/runitcmd/cmd/runitcmd

To install for your system

    sudo -E GOPATH=/ go get -u github.com/sigmonsays/runitcmd/cmd/runitcmd


# configuration

configuration is optional and sane defaults will be used when no configuration is present

the system wide configuration file is /etc/runitcmd.yaml. The user configuration overrides the system
wide configuration.

example user configuration:

      sudo: true
      logging:
        directory: /opt/logs
        number: 10
        max_size: 52428800

# usage

      % runitcmd -h
      NAME:
         runitcmd - manage runit services

      USAGE:
         runitcmd [global options] command [command options] [arguments...]

      VERSION:
         0.0.1

      COMMANDS:
           list, ls    list services
           create      Create services
           setup       setup a service
           apply       update service
           status, st  get service status
           import      import service
           export      export service
           delete      delete a service
           activate    activate a service
           deactivate  deactivate a service
           enable      enable a service
           disable     disable a service
           reset       reset a service
           up          up a service
           down        down a service
           pause       pause a service
           cont        cont a service
           hup         hup a service
           alarm       alarm a service
           interrupt   interrupt a service
           quit        quit a service
           usr1        usr1 a service
           usr2        usr2 a service
           term        term a service
           kill        kill a service
           start       start a service
           stop        stop a service
           reload      reload a service
           restart     restart a service
           shutdown    shutdown a service
           help, h     Shows a list of commands or help for one command

      GLOBAL OPTIONS:
         --config value, -c value  override configuration file
         --level value, -l value   change log level (default: "WARN")
         --service-dir value       change service dir
         --active-dir value        change active service dir
         --help, -h                show help
         --version, -v             print the version

# Examples

Setup a service (and start it)

    runitcmd setup sleep --run 'sleep 3600'

Stop/start service

    runitcmd stop sleep
    runitcmd start sleep

Delete service

    runitcmd delete sleep

`EOF`

