FROM nodered/node-red:latest-22

# 1. Switch to root to install system dependencies
USER root

# 2. Install build tools and python (needed for certain node compilations)
RUN apk add --no-cache build-base python3 

# 3. Set the environment variable for your custom config location
ENV CONFIG_FILE_PATH=/data/mixer_config.json

# 4. Set the working directory to the standard Node-RED app folder
WORKDIR /usr/src/node-red

# 5. Copy the Logic files to the internal App folder (NOT /data)
# This ensures these files are part of the image layers and not hidden by volumes
COPY --chown=node-red:node-red ./src/package.json ./package.json
COPY --chown=node-red:node-red ./src/flows.json ./flows.json

# 6. Switch back to the node-red user to install packages
USER node-red

# 7. Install your specific palettes into the App folder
RUN npm install @flowfuse/node-red-dashboard \
    node-red-contrib-osc 

# 8. Copy the TEMPLATE config to /data
# Note: This will only be used if the volume is empty. 
# Once a volume exists at the church, this line is effectively ignored by Podman.
COPY --chown=node-red:node-red ./src/mixer_config.json /data/mixer_config.json