# Import the installable/sh image
FROM installable/sh AS installable AS run

# Install Claude Code using the RUN binary
FROM ubuntu:latest AS claude
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
COPY --from=installable / /
RUN RUN "https://claude.ai/install.sh"

# Create a clean final image with Claude Code installed
FROM ubuntu:latest
COPY --from=claude / /
ENV PATH="/root/.local/bin:$PATH"
ENTRYPOINT ["claude"]
CMD ["--help"]
