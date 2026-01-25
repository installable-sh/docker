FROM ghcr.io/scaffoldly/run AS run

FROM ubuntu:latest
COPY --from=run / /
CMD ["RUN", "https://raw.githubusercontent.com/scaffoldly/run/main/examples/hello-world.sh"]
