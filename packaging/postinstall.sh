#!/bin/sh

systemctl daemon-reload

systemctl enable updates-exporter.service
