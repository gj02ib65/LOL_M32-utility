This project was born out of a specific need for high-quality livestream audio without doubling the workload of the sound crew. By using a second Midas/Behringer console dedicated to the stream, we gained independent control over EQ, compression, and levels—but found that managing the "mutes" on two separate desks during a live event was prone to error.

This application bridges that gap, allowing a single engineer to handle the Front of House (FOH) mix while the system automatically handles the tedious task of mirroring mutes and channel data to the Livestream console.

---

# Midas/X32 Dual Mixer Sync & Stream Manager

## 📖 The Problem & Purpose

When running a livestream alongside a live event, the audio needs are vastly different. In the room, you mix for acoustics; for the stream, you mix for speakers. While a dedicated second mixer is the ideal solution, it usually requires a second engineer just to keep the microphones in sync.

**This application allows one audio engineer to effectively manage both mixers.** Once the livestream levels and EQs are set, the software takes over the "chore" of muting and unmuting channels. When you mute a microphone on the FOH console, it is instantly muted on the Livestream console, ensuring the stream never hears "hot" mics that aren't in use.

---

## 🚀 Key Features

* **Synchronized Muting:** Real-time replication of "Mute" commands for all 32 channels and 8 DCAs.
* **Metadata Sync:** Automatically mirrors Channel Names and DCA Names from the Primary console to the Replica.
* **DCA Management:** Replicates DCA assignments and DCA master mutes.
* **Safe Channels:** A keyword-based exclusion system. Any channel containing a "Safe" substring (e.g., "Music," "Ambient") will ignore incoming sync commands, allowing independent muting for the stream.
* **Active Watchdog:** Visual feedback on the dashboard showing "Online," "Offline," or "Simulated" status for both consoles.
* **Simulation Mode:** Process sync logic internally for testing without actually sending commands to the hardware.

---

## 🛠 Project Architecture

The system is built on a high-performance **UDP/OSC Routing Engine** designed to handle the high traffic of a digital console without adding latency or network jitter.

### Core Components:

1. **Network Router:** A low-latency "Traffic Cop" that identifies incoming OSC packets. It intentionally filters out fader movements and meters to prevent network congestion, focusing only on mutes and configuration data.
2. **Sync Logic Engine:** Compares incoming commands against the **Safe List**. If a channel is "Safe," the sync command is blocked.
3. **Persistence Layer:** All settings (IP addresses, Safe keywords, and Mode states) are saved to a local `mixer_config.json` file so the system is ready to go immediately upon power-up.
4. **Startup Synchronizer:** On boot, the application performs a "Handshake" query to the Primary console to pull current names and mute states, ensuring both desks start in perfect alignment.

---

## 🚢 Installation & Deployment

For end-users who simply want to run the utility without modifying the source code, the easiest method is to use Docker/Podman with the pre-built image.

1. **Install Docker or Podman** on your server/computer.
2. **Create a `docker-compose.yml` file** with the following content:

```yaml
services:
  m32-sync:
    image: ghcr.io/gj02ib65/lol_m32-utility:latest
    container_name: m32-sync-utility
    restart: always
    network_mode: "host" 
    ports:
      - "1880:1880"
      - "10023:10023/udp"
    volumes:
      - m32_church_data:/data 
    environment:
      - TZ=America/Chicago

  # Optional: Watchtower automatically applies monthly security updates
  watchtower:
    image: containrrr/watchtower
    container_name: watchtower
    restart: always
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --schedule "0 0 3 * * 1" --cleanup m32-sync-utility

volumes:
  m32_church_data: 
```

3. **Run the application**: Open your terminal in the same folder as the file and run `docker compose up -d` (or `podman compose up -d`).
4. **Access the Dashboard**: Open your web browser and navigate to `http://localhost:1880/dashboard/mixer`.

*(Developers: See the developer section below for instructions on building from source.)*

---

## 💻 Developer Instructions (Building from Source)

If you wish to modify the Node-RED flow (`src/flows.json`), you will want to build the container locally. This project utilizes [Task](https://taskfile.dev/) (also known as `go-task`) to run local build and sync scripts.

### 1. Install Task (`go-task`)
Depending on your platform, install the Task runner:
- **Windows (Scoop):** `scoop install task`
- **macOS (Homebrew):** `brew install go-task`
- **Linux:** [See installation guide](https://taskfile.dev/installation/)

### 2. Available Commands
You can run the following commands from the root directory:
- `task build` -> Rebuilds the container locally using the source files in `/src`.
- `task up` / `task down` -> Start or stop the local development container.
- `task status` -> Check container health and the currently loaded config.
- `task logs` -> Follow the live container logs.
- `task pull` -> Syncs the live `flows.json` and `package.json` from your running container back into your local `/src` Git repository.

---

## ⚙️ Usage Configuration

1. **Primary IP:** The address of your FOH Console.
2. **Replica IP:** The address of your Livestream Console.
3. **Safe Substrings:** Comma-separated list of names to protect (e.g., `Host, Video, BGM`). If a channel name contains any of these words, its mute status will **not** be synchronized.
4. **Simulation Mode:** Toggle this to "On" to test your logic or Safe List without affecting the Livestream console.

---

## ⚠️ Requirements

* **Hardware:** Two Behringer X32 or Midas M32 series consoles on the same network.
* **Platform:** Node-RED with `node-red-contrib-osc` and `node-red-dashboard` 2.0.
* **Port:** Consoles must be reachable on UDP Port `10023`.

---

## 🧠 Logic Flow: The Journey of a Mute Command

To ensure stability and performance, every command follows a strict "Filter and Verify" path. Here is what happens when you press a **Mute** button on the Front of House (Primary) console:

### 1. The Trigger (FOH Console)

The FOH console sends an OSC packet over UDP (Port 10023) to Node-RED. For example, muting Channel 5 sends:
`topic: /ch/05/mix/on`, `payload: 0`

### 2. The Traffic Cop (Network Router)

The message enters the **Network Router**. This node performs three critical checks:

* **Origin Check:** Does the source IP match the **Primary IP**? (If no, drop it).
* **Fader Filter:** Is this a fader movement or meter data? (If yes, drop it to save bandwidth).
* **Routing:** Proof-of-life packets are sent to the **Watchdog**, while valid Sync commands (Mutes/Names) are passed to the **Sync Logic Engine**.

### 3. The Brain (Sync Logic Engine)

This is where the "Safe List" and "Simulation" rules are applied:

* **Name Lookup:** The engine checks the current name of Channel 5 (stored in Global memory).
* **The Safe Check:** If the name is "Host Mic" and your Safe List contains "Host," the logic **blocks** the message. The FOH mic mutes, but the Livestream mic stays on.
* **Simulation Check:** If **Simulation Mode** is active, the command is diverted to the dashboard logs only and never hits the network.

### 4. The Handshake (Replica Console)

If the command passes all checks, it is encoded back into a raw OSC packet and sent to the **Replica IP**. The Livestream console receives the command and instantly updates its mute state to match the FOH desk.

---

## 🛠 Project Data Hierarchy

| Priority | Level | Description |
| --- | --- | --- |
| **1** | **Online Proof** | Real-time traffic from a mixer always overrides any other status. |
| **2** | **Simulation Mode** | If hardware is missing and Sim Mode is ON, the dashboard shows "Simulated." |
| **3** | **Offline** | If no traffic is heard for 10 seconds and Sim Mode is OFF, the desk is marked "Offline." |

---

## 📋 Summary of Replicated Data

Beyond simple mutes, the system ensures your entire console "profile" is mirrored where it matters most:

* **Channel Names:** `/ch/xx/config/name`
* **Channel Mutes:** `/ch/xx/mix/on`
* **DCA Names:** `/dca/x/config/name`
* **DCA Mutes:** `/dca/x/on`
* **DCA Assignments:** `/ch/xx/grp/dca` (Ensures your channel-to-DCA mapping stays identical).

---
