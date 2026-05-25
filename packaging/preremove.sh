#!/bin/sh

systemctl stop updates-exporter.service || true
systemctl disable updates-exporter.service || true
