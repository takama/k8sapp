FROM scratch

COPY certs /etc/ssl/
COPY bin/linux-amd64/k8sapp /

CMD ["/k8sapp"]
