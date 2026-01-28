# Import the RUN image
FROM installable/sh AS run

# Use any base image with the utilities your script needs
FROM ubuntu:latest

# Add the RUN and INSTALL binaries to the image
COPY --from=run / /

# Fetch and execute a shell script from a URL
CMD ["RUN", "https://raw.githubusercontent.com/installable-sh/docker/main/examples/hello-world/hello-world.sh"]
