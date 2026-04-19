FROM nodered/node-red:latest-22

USER root
RUN apk add --no-cache build-base python3 

# 1. Prepare your logic folder
RUN mkdir -p /opt/m32-logic && chown node-red:node-red /opt/m32-logic

# 2. Copy and Install palettes into your logic folder
# We use a temporary WORKDIR just for the installation
WORKDIR /opt/m32-logic
COPY --chown=node-red:node-red ./src/package.json ./package.json
COPY --chown=node-red:node-red ./src/flows.json ./flows.json

USER node-red
RUN npm install --omit=dev

# 3. CRITICAL: Switch WORKDIR back to the base image default
# This allows the container to find './entrypoint.sh' on startup
WORKDIR /usr/src/node-red

# 4. Set paths for persistent config
ENV CONFIG_FILE_PATH=/data/mixer_config.json