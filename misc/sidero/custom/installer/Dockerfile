FROM scratch AS customization
COPY --from=ghcr.io/charlie-haley/kernel:v0.9.0-alpha.0-25-g49f3c17-dirty /lib/modules /lib/modules

FROM ghcr.io/talos-systems/installer:latest
COPY --from=ghcr.io/charlie-haley/kernel:v0.9.0-alpha.0-25-g49f3c17-dirty /boot/vmlinuz /usr/install/${TARGETARCH}/vmlinuz
