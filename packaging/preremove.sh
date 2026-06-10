#!/bin/sh

systemctl stop updates_exporter.service || true
systemctl disable updates_exporter.service || true
