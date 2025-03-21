# api

## 创建表

```sql
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    image TEXT NOT NULL, 
    description TEXT NOT NULL,
    url TEXT NOT NULL
);
```