const fs = require('fs');
let data = JSON.parse(fs.readFileSync('src/flows.json', 'utf8'));

let count = 0;
data.forEach(n => {
    if (n.type === 'function' && n.func) {
        // Look for literal "\n" strings that were accidentally inserted
        // Note: In Javascript, the string '\\n' represents the literal slash-n sequence
        if (n.func.includes('\\n')) {
            count++;
            n.func = n.func.replace(/\\n/g, '\n');
        }
    }
});

fs.writeFileSync('src/flows.json', JSON.stringify(data, null, 4));
console.log(`Cleaned up escaped newlines in ${count} function nodes.`);
