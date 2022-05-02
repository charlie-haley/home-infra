FROM --platform=amd64 ghcr.io/charlie-haley/kernel:v0.10.0-alpha.0-5-g54aa902-dirty AS kernel
FROM --platform=amd64 ghcr.io/charlie-haley/nvidia-container-runtime:v0.10.0-alpha.0-5-g54aa902-dirty AS pkg-nvidia-container-runtime

FROM scratch AS customization
COPY --from=kernel /lib/modules /lib/modules
COPY --from=pkg-nvidia-container-runtime /bin/nvidia-container-runtime /bin/nvidia-container-runtime

FROM ghcr.io/talos-systems/installer:v1.0.3

COPY --from=ghcr.io/charlie-haley/kernel:v0.10.0-alpha.0-5-g54aa902-dirty /boot/vmlinuz /usr/install/${TARGETARCH}/vmlinuz
