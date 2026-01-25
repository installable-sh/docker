# Import the RUN image
FROM scaffoldly/run AS run

# Use any base image with the utilities your script needs
FROM ubuntu:latest

# Add the RUN binary to the image
COPY --from=run / /

# Fetch and execute a shell script from a URL
CMD ["RUN", "https://raw.githubusercontent.com/scaffoldly/run/main/examples/hello-world/hello-world.sh"]
