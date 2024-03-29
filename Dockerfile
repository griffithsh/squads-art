FROM fedora:latest

RUN dnf install -y \
    go \
    gimp \
    tiled \
    xorg-x11-server-Xvfb

WORKDIR /mnt

CMD xvfb-run ./everything.sh
