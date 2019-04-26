#!/bin/bash

set -eou pipefail

package_name="zabbix-agent-extension-sentry"

rm -rf *.tar.xz
makepkg -Cod; PKGVER=$(cd $(pwd)/src/$package_name/ && make ver) makepkg -esd
