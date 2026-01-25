# Import the RUN image
FROM scaffoldly/run AS run

# Use a base image with bash installed
FROM ubuntu:latest
COPY --from=run / /
CMD ["RUN", "https://raw.githubusercontent.com/scaffoldly/run/main/examples/bash/in-bash.sh"]
