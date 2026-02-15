FROM ghcr.io/installable-sh/run:v0 AS run
FROM ghcr.io/installable-sh/install:v0 AS install

FROM scratch AS combined
COPY --from=run /usr/local/bin/RUN /usr/local/bin/RUN
COPY --from=install /usr/local/bin/INSTALL /usr/local/bin/INSTALL
RUN ["/usr/local/bin/RUN", "--help"]
RUN ["/usr/local/bin/INSTALL", "--help"]

FROM scratch
COPY --from=combined / /
ENTRYPOINT ["RUN"]
CMD ["--help"]
