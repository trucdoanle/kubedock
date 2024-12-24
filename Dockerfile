FROM us-east4-docker.pkg.dev/sym-prod-mr-tools-01/jenkins-docker-us-east4/ubuntu:noble.Production-146-8cc6c3f

COPY kubedock /usr/local/bin

ENTRYPOINT ["/usr/local/bin/kubedock"]
CMD [ "server" ]
