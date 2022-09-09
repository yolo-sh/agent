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

# -- Set hostname

sudo hostnamectl set-hostname "${ENV_NAME_SLUG}"

# -- System dependencies

# log "Installing system and Docker dependencies"

sudo apt-get --assume-yes --quiet --quiet install wget git vim apt-transport-https ca-certificates gnupg lsb-release

# -- Yolo volume configuration

# log "Locating Yolo volume"

# YOLO_VOL_MOUNTPOINT="/yolo"
# YOLO_VOL_LABEL="yolo-volume"
# DOCKER_DATA_DIR="${YOLO_VOL_MOUNTPOINT}/docker"

# # Find the volume not mounted (mountpoint == null) 
# # and with no partition (children == null)
# YOLO_VOL=$(lsblk --json | jq '.blockdevices[] | select(.mountpoint == null and .children == null) | .name' --raw-output)

# if [[ "${YOLO_VOL}" = "" ]]; then
#   echo "Yolo volume not found"
#   exit 1
# elif [[ "${YOLO_VOL}" = *" "* ]]; then # eg: sda1 sda2
#   echo "Multiple volumes match the Yolo one"
#   exit 1
# fi

# YOLO_VOL="/dev/${YOLO_VOL}"

# log "Yolo volume found ${YOLO_VOL}"

# # If the output shows simply data, there is no file system on the device
# # Example:
# # [ec2-user ~]$ file --special-files /dev/xvdf
# # /dev/xvdf: data
# YOLO_VOL_IS_FORMATTED=$([[ $(sudo file --special-files "${YOLO_VOL}") = "${YOLO_VOL}: data" ]] && echo "false" || echo "true")

# if [[ "${YOLO_VOL_IS_FORMATTED}" = "false" ]]; then
#   log "Yolo volume not formatted. Formatting now..."
#   sudo mkfs.ext4 "${YOLO_VOL}"
# fi

# log "Backuping /etc/fstab to /etc/fstab.orig"
# sudo cp /etc/fstab /etc/fstab.orig

# log "Labeling Yolo volume to ${YOLO_VOL_LABEL}"
# sudo e2label "${YOLO_VOL}" "${YOLO_VOL_LABEL}"

# log "Creating ${YOLO_VOL_MOUNTPOINT} mountpoint"
# sudo mkdir --parents "${YOLO_VOL_MOUNTPOINT}"

# log "Adding Yolo volume to fstab"
# echo "LABEL=${YOLO_VOL_LABEL}  ${YOLO_VOL_MOUNTPOINT}  ext4  defaults,nofail  0  2" | sudo tee --append /etc/fstab > /dev/null

# log "Mounting all devices"
# sudo mount --all

# if [[ "${YOLO_VOL_IS_FORMATTED}" = "true" ]]; then
#   log "Yolo volume already formatted before mounting. Making sure the filesystem size match the attached volume..."
#   sudo resize2fs "${YOLO_VOL}"
# fi

# sudo mkdir --parents "${DOCKER_DATA_DIR}"
# sudo chown --recursive yolo:yolo "${YOLO_VOL_MOUNTPOINT}"

# -- Packages configuration

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

# log "Configuring Docker"

# sudo mkdir /etc/systemd/system/docker.service.d
# sudo touch /etc/systemd/system/docker.service.d/override.conf

# sudo tee /etc/systemd/system/docker.service.d/override.conf > /dev/null << EOF
# [Service]
# ExecStart=
# ExecStart=/usr/bin/dockerd --host fd:// --containerd=/run/containerd/containerd.sock --data-root="${DOCKER_DATA_DIR}"
# EOF

# sudo systemctl daemon-reload
# sudo systemctl restart docker.service

# -- Run as "yolo"

log "Configuring workspace for user \"yolo\""

sudo --set-home --login --user yolo -- env \
	GITHUB_USER_EMAIL="${GITHUB_USER_EMAIL}" \
	USER_FULL_NAME="${USER_FULL_NAME}" \
bash << 'EOF'

mkdir --parents .vscode-server

if [[ ! -f ".ssh/yolo_github" ]]; then
	ssh-keygen -t ed25519 -C "${GITHUB_USER_EMAIL}" -f .ssh/yolo_github -q -N ""
fi

chmod 644 .ssh/yolo_github.pub
chmod 600 .ssh/yolo_github

if ! grep --silent --fixed-strings "IdentityFile ~/.ssh/yolo_github" .ssh/config; then
	rm --force .ssh/config
  echo "Host github.com" >> .ssh/config
	echo "  User git" >> .ssh/config
	echo "  Hostname github.com" >> .ssh/config
	echo "  PreferredAuthentications publickey" >> .ssh/config
	echo "  IdentityFile ~/.ssh/yolo_github" >> .ssh/config
fi

chmod 600 .ssh/config

if ! grep --silent --fixed-strings "github.com" .ssh/known_hosts; then
  ssh-keyscan github.com >> .ssh/known_hosts
fi

GIT_GPG_KEY_COUNT="$(gpg --list-signatures --with-colons | grep 'sig' | grep "${GITHUB_USER_EMAIL}" | wc -l)"

if [[ $GIT_GPG_KEY_COUNT -eq 0 ]]; then
	gpg --quiet --batch --gen-key << EOF2
%no-protection
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: ${USER_FULL_NAME}
Name-Email: ${GITHUB_USER_EMAIL}
Expire-Date: 0
EOF2
fi

GIT_GPG_KEY_ID="$(gpg --list-signatures --with-colons | grep 'sig' | grep "${GITHUB_USER_EMAIL}" | head --lines 1 | cut --delimiter ':' --fields 5)"

if [[ ! -f ".gnupg/yolo_github_gpg_public.pgp" ]]; then
	GIT_GPG_PUBLIC_KEY="$(gpg --armor --export "${GIT_GPG_KEY_ID}")"

	echo "${GIT_GPG_PUBLIC_KEY}" >> .gnupg/yolo_github_gpg_public.pgp
fi

chmod 644 .gnupg/yolo_github_gpg_public.pgp

if [[ ! -f ".gnupg/yolo_github_gpg_private.pgp" ]]; then
	GIT_GPG_PRIVATE_KEY="$(gpg --armor --export-secret-keys "${GIT_GPG_KEY_ID}")"

	echo "${GIT_GPG_PRIVATE_KEY}" >> .gnupg/yolo_github_gpg_private.pgp
fi

chmod 600 .gnupg/yolo_github_gpg_private.pgp

git config --global pull.rebase false

git config --global user.name "${USER_FULL_NAME}"
git config --global user.email "${GITHUB_USER_EMAIL}"

git config --global user.signingkey "${GIT_GPG_KEY_ID}"
git config --global commit.gpgsign true

EOF
