#!/bin/bash
# Yolo instance init
set -euo pipefail

log () {
  echo -e "${1}" >&2
}

# Remove "debconf: unable to initialize frontend: Dialog" warnings
echo 'debconf debconf/frontend select Noninteractive' | sudo tee debconf-set-selections > /dev/null

handleExit () {
  EXIT_CODE=$?
  exit "${EXIT_CODE}"
}

trap "handleExit" EXIT

# -- System configuration

# Lookup instance architecture for Sysbox
INSTANCE_ARCH=""
case $(uname -m) in
  i386)       INSTANCE_ARCH="386" ;;
  i686)       INSTANCE_ARCH="386" ;;
  x86_64)     INSTANCE_ARCH="amd64" ;;
  arm)        dpkg --print-architecture | grep -q "arm64" && INSTANCE_ARCH="arm64" || INSTANCE_ARCH="armv6" ;;
  aarch64_be) INSTANCE_ARCH="arm64" ;;
  aarch64)    INSTANCE_ARCH="arm64" ;;
  armv8b)     INSTANCE_ARCH="arm64" ;;
  armv8l)     INSTANCE_ARCH="arm64" ;;
esac

# -- Packages dependencies

log "Installing Docker / Sysbox dependencies"

sudo apt-get --assume-yes --quiet --quiet install ca-certificates curl gnupg jq lsb-release wget

# Docker
log "Installing Docker"

if [[ ! -f "/usr/share/keyrings/docker-archive-keyring.gpg" ]]; then
  curl --fail --silent --show-error --location https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor --output /usr/share/keyrings/docker-archive-keyring.gpg
fi

if [[ ! -f "/etc/apt/sources.list.d/docker.list" ]]; then
	echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release --codename --short) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
fi

sudo apt-get --assume-yes --quiet --quiet update
sudo apt-get --assume-yes --quiet --quiet remove docker docker-engine docker.io containerd runc
sudo apt-get --assume-yes --quiet --quiet install docker-ce docker-ce-cli containerd.io

# Sysbox
log "Installing Sysbox"

if [[ ! -f "/tmp/sysbox.deb" ]]; then
  wget "https://downloads.nestybox.com/sysbox/releases/v0.5.2/sysbox-ce_0.5.2-0.linux_${INSTANCE_ARCH}.deb" -O /tmp/sysbox.deb
fi

sudo apt-get --assume-yes --quiet --quiet install /tmp/sysbox.deb
rm -rf /tmp/sysbox.deb