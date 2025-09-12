# Prometheus Updates Exporter

This is a Prometheus exporter that exposes system package update statistics:

* *pending_updates*: total number of pending updates
* *pending_security_updates*: pending security updates
* *reboot_required*: is reboot required due (due to package updates)

RHEL derivatives and Ubuntu are tested and supported. Other dnf and apt-based distros probably work fine.

# Installation

The easiest way is to install a prebuilt RPM package:

```
$ wget https://puppeteers-public.hel1.your-objectstorage.com/updates_exporter-0.0.7-1.x86_64.rpm
$ rpm -i updates_exporter-0.0.7-1.x86_64.rpm
$ systemctl enable updates-exporter
$ systemctl start updates-exporter
```

Alternatively build yourself or get a precompiled executable from the GitHub
Releases page. You can use
[podman-builder](https://github.com/Puppet-Finland/podman-builder) to build
your own package.

# Usage

Updates Exporter listens on port 9091. Just scrape /metrics and you're good to go.
