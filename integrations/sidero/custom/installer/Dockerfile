FROM scratch AS customization
COPY --from=ghcr.io/charlie-haley/kernel:v0.10.0-alpha.0-5-g54aa902-dirty /lib/modules /lib/modules

FROM ghcr.io/talos-systems/installer:latest
COPY --from=ghcr.io/charlie-haley/kernel:v0.10.0-alpha.0-5-g54aa902-dirty /boot/vmlinuz /usr/install/${TARGETARCH}/vmlinuz
