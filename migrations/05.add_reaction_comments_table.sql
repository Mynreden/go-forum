CREATE TABLE IF NOT EXISTS commentsReactions (
    id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    comment_id  INTEGER NOT NULL,
    user_id     INTEGER NOT NULL,
    reaction   INTEGER NOT NULL CHECK(reaction IN (0, 1)),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);