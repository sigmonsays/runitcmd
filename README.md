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

    GOPATH=/ go get -u github.com/sigmonsays/runitcmd/cmd/runitcmd

# configuration

configuration is optional and sane defaults will be used when no configuration is present

the system wide configuration file is /etc/runitcmd.yaml. The user configuration overrides the system
wide configuration.

example user configuration:

      sudo: true



