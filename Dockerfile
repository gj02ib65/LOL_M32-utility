FROM nodered/node-red:latest

# 1. Switch to root to install system dependencies
USER root

# 2. Install ALSA headers and build tools
# 'alsa-lib-dev': Contains alsa/asoundlib.h required by the midi module
# 'build-base': Contains gcc/g++/make required to compile the C++ code
RUN apk add --no-cache build-base python3 tcpdump doas
# RUN adduser -D -s /bin/sh node-red
RUN mkdir -p /etc/doas.d && \
    echo "permit nopass node-red as root" > /etc/doas.d/doas.conf
# 3. Switch back to the node-red user
USER node-red

# 4. Install the Node-RED packages
RUN npm install @flowfuse/node-red-dashboard \
                node-red-contrib-osc 
                # node-red-contrib-midi
                
COPY --chown=node-red:node-red ./src/flows.json /data/flows.json
COPY --chown=node-red:node-red ./src/mixer_config.json /data/mixer_config.json