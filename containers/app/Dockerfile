FROM nodered/node-red:latest-22

# 1. Switch to root for system-level installs
USER root

# Restoring your critical tools: 
# build-base & python3: Needed to compile OSC/MIDI C++ components
# alsa-lib-dev: Critical for MIDI/OSC communication stability
# tcpdump: For your network troubleshooting
RUN apk add --no-cache alsa-lib-dev build-base python3 tcpdump

# 2. Create the "Baked-in" logic folder and fix permissions immediately
RUN mkdir -p /opt/m32-logic && chown -R node-red:node-red /opt/m32-logic

# 3. Move to the logic folder to perform the installation
WORKDIR /opt/m32-logic

# 4. Copy your current manifest and flows
COPY --chown=node-red:node-red ./src/package.json ./package.json
COPY --chown=node-red:node-red ./src/flows.json ./flows.json

# 5. Switch to node-red user for the npm install
USER node-red

# We force the installation of the specific palettes. 
# This ensures they are physically present in /opt/m32-logic/node_modules
RUN npm install --omit=dev @flowfuse/node-red-dashboard \
    node-red-contrib-osc

# 6. CRITICAL: Switch WORKDIR back to the default
# This ensures the container can find its internal startup scripts (entrypoint.sh)
WORKDIR /usr/src/node-red

# 7. Final environment setup
ENV CONFIG_FILE_PATH=/data/mixer_config.json