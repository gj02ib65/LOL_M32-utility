FROM nodered/node-red:latest-22

USER root
# Install system tools
RUN apk add --no-cache alsa-lib-dev build-base python3 tcpdump

# 1. Setup the baked-in folder
RUN mkdir -p /opt/m32-logic && chown -R node-red:node-red /opt/m32-logic
WORKDIR /opt/m32-logic

# 2. Copy logic files
COPY --chown=node-red:node-red ./src/package.json ./package.json
COPY --chown=node-red:node-red ./src/flows.json ./flows.json

USER node-red

# 3. Install palettes locally in /opt/m32-logic
RUN npm install --omit=dev

# 4. CRITICAL: Tell Node.js to ALWAYS look here for modules
ENV NODE_PATH=/opt/m32-logic/node_modules
ENV CONFIG_FILE_PATH=/data/mixer_config.json

# 5. Reset WORKDIR for the entrypoint script
WORKDIR /usr/src/node-red