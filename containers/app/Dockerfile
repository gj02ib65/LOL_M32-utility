FROM nodered/node-red:latest-22

USER root
# Install build tools (needed for OSC/MIDI nodes)
RUN apk add --no-cache build-base python3 

# 1. Create a dedicated folder for your BAKED-IN logic
RUN mkdir -p /opt/m32-logic && chown node-red:node-red /opt/m32-logic
WORKDIR /opt/m32-logic

# 2. Copy your logic files here (NOT to /usr/src/node-red)
COPY --chown=node-red:node-red ./src/package.json ./package.json
COPY --chown=node-red:node-red ./src/flows.json ./flows.json

USER node-red

# 3. Install your palettes into this dedicated folder
RUN npm install --omit=dev @flowfuse/node-red-dashboard \
    node-red-contrib-osc

# 4. Set the path for your persistent mixer config
ENV CONFIG_FILE_PATH=/data/mixer_config.json