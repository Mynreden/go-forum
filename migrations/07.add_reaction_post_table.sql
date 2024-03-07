CREATE TABLE IF NOT EXISTS postsReactions (
    id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    post_id     INTEGER NOT NULL, 
    user_id     INTEGER NOT NULL,
    reaction    INTEGER NOT NULL CHECK(reaction IN (0, 1)),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)