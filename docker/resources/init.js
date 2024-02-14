db = db.getSiblingDB('dummy');
db.c1.insert({ key: 'value1', number: 12, bool: true });
db.c1.insert({ key: 'value2', number: 13, bool: true });
db.c1.insert({ key: 'value3', number: 14 });
db.c1.insert({ key: 'value4', number: 15, bool: true });

db.c2.insert({ complex: { name: "homer", array: [1, 2, 3, 4] }, bool: false });
db.c2.insert({ complex: { name: "marge", array: [1, 2, 3] }, bool: true });
db.c2.insert({ complex: { name: "maggy", array: [1, 2, 3], hobbies: { saxophon: true, skating: false } }, bool: true });
