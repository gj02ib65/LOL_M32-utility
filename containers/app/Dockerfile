FROM nodered/node-red:latest-22

USER root
# Install system tools (provides gcc, python3, make for native npm modules)
RUN apk add --no-cache alsa-lib-dev build-base python3 tcpdump

# 1. Create the m32-logic directory — holds palettes and settings.
#    flows.json is NOT stored here. It is sourced from the bind-mounted
#    directory at /opt/m32-flows-src/ and copied to /data/ at startup.
RUN mkdir -p /opt/m32-logic && chown -R node-red:node-red /opt/m32-logic

# 2. Copy the Node-RED settings file.
#    This tells Node-RED's node scanner to look in /opt/m32-logic/node_modules
#    for installed palette nodes (via the nodesDir setting).
COPY --chown=node-red:node-red ./src/settings.js /opt/m32-logic/settings.js
# Default config — copied to /data/mixer_config.json only on first run (never overwrites).
COPY --chown=node-red:node-red ./src/mixer_config.json /opt/m32-logic/mixer_config.default.json

# 3. Copy package.json to the prefix location and install palettes.
#    npm install --prefix reads package.json FROM the prefix directory.
#    Run npm from /tmp to avoid Node-RED's bundled npm (which is missing node-gyp).
COPY --chown=node-red:node-red ./src/package.json /opt/m32-logic/package.json
WORKDIR /tmp/palette-install
RUN npm install --prefix /opt/m32-logic --omit=dev && \
    chown -R node-red:node-red /opt/m32-logic

# 4. Install the startup script.
#    This copies flows.json from the bind-mounted source directory into /data/
#    before Node-RED starts, avoiding the EBUSY rename error that occurs when
#    Node-RED tries to atomically save a directly bind-mounted file.
COPY ./containers/app/entrypoint.sh /m32-entrypoint.sh
RUN chmod +x /m32-entrypoint.sh

USER node-red

# 5. Runtime environment.
#    FLOWS points to /data/flows.json — the writable copy created by the entrypoint.
ENV CONFIG_FILE_PATH=/data/mixer_config.json
ENV FLOWS=/data/flows.json

ENTRYPOINT ["/m32-entrypoint.sh"]
CMD []

WORKDIR /usr/src/node-red