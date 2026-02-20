const fs = require('fs');
let data = JSON.parse(fs.readFileSync('src/flows.json'));

data.forEach(n => {
  if (n.id === 'cb60e1c7730b4854') { // Get IPs and Xremote Heartbeat
    if (!n.func.includes('global.get("system_active")')) {
      n.func = 'if (global.get("system_active") === false) {\\n    node.status({fill:"yellow", shape:"ring", text:"Standby"});\\n    return null;\\n}\\n\\n' + n.func;
    }
  }
  if (n.id === '6b36dad599913889') { // Network Router
    if (!n.func.includes('global.get("system_active")')) {
      n.func = 'if (global.get("system_active") === false) return [null, null, null];\\n\\n' + n.func;
    }
  }
  if (n.id === '53180069e840ae5b') { // Query Primary Desk
     if (!n.func.includes('global.get("system_active")')) {
      n.func = 'if (global.get("system_active") === false) return null;\\n\\n' + n.func;
    }
  }
});

fs.writeFileSync('src/flows.json', JSON.stringify(data, null, 4));
console.log('flows.json updated successfully.');
