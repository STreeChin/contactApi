//db = db.getSiblingDB('contact'); // creat DB: contact

db.createUser(
    {
        user: "user",
        pwd: "user",
        roles: [
            { 
                role: "readWrite", 
                db: "contact"
            }
        ]
    }
);

db.createCollection("contactInfo"); 

db.createCollection("apiKey"); 

db.apiKey.insertOne( 
    { 
        contactid: BinData(0,"0D60XbKwTnK7OOsG2pMR43LobAnUjJ3tN9E+ARIfpHCCm3yGjY0kTpZLHFvjt2Hu"), 
        apikey: BinData(0,"yaKQtPiK4SekQEJv7jx/5OKOs6wWD9azfM1ba0wKnIgxiPoL6xn1jGjOAFfd1NEM")
    } 
);
db.contactInfo.insertOne( 
    { 
        contactid: "test init id", 
        email: "init@test.com"
    } 
);

