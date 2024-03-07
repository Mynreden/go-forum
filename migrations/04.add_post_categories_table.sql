CREATE TABLE IF NOT EXISTS categories (
    id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    category_name   VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS PostCategories (
  post_id INTEGER,
  category_name INTEGER,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (category_name) REFERENCES categories(category_name) ON DELETE CASCADE
);
