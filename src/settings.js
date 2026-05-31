// Node-RED settings for M32 Utility
// This is a minimal settings file. Its primary job is to tell Node-RED
// where to find the palette nodes that are installed at build time.
module.exports = {

    // Tell Node-RED to scan /opt/m32-logic/node_modules for installed palette nodes.
    // Node-RED looks for a node_modules sub-directory inside each path listed here.
    // npm installs our palettes to /opt/m32-logic/node_modules/ via --prefix /opt/m32-logic.
    nodesDir: ['/opt/m32-logic'],

    // All other settings use Node-RED's built-in defaults.
}
