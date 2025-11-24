CREATE TABLE IF NOT EXISTS Users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    team TEXT NOT NULL REFERENCES Team(name)
); 
