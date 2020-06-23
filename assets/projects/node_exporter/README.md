# Getting Started

The `config.json` is a configuration file for `gorpm` utility.

In the `config.json`, under `files` section, we define the files to install:

* `src/node_exporter-0.18.1.linux-amd64/node_exporter`: main executable that goes to `/usr/local/bin/`
* `etc/profile.d/node-exporter.sh`: is the `/etc/profile.d/` file with environment variables
* `etc/sysconfig/node-exporter.conf`: configuration file
* `systemd/node-exporter.service`: systemd unit file

Additionally, there are references to the following scripts. The scripts are
being invoked at different phases of RPM lifecycle, e.g. install, removal, etc.

* `scripts/pre_install.sh`
* `scripts/post_install.sh`
* `scripts/pre_remove.sh`
* `scripts/post_remove.sh`
* `scripts/verify.sh`
* `scripts/cleanup.sh`: runs during `rpmbuild`
