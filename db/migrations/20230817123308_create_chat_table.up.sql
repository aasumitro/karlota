-- Conversations Table
CREATE TABLE IF NOT EXISTS conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    type text DEFAULT 'private' CHECK(type IN ('private', 'group') ),
    name VARCHAR(100), -- for group user - admin can add name
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Participants Table
CREATE TABLE IF NOT EXISTS participants(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    conversation_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    -- if this for group conversation then we will have 1 or more admin, and also 1 or more member
    -- group will be automatically created, when user create new conversation and the recipient more than 2 then group will be craeted
    role text DEFAULT 'none' CHECK(role IN ('none', 'admin', 'member')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL -- this used only for private conversation for group different case
);

-- Messages Table
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    conversation_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    type text DEFAULT 'text' CHECK(type IN ('text', 'image', 'file', 'voice', 'voice_call', 'video_call')),
    content text NOT NULL,
    read_at DATETIME DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance improvement
CREATE INDEX idx_participants_conversation_id ON participants(conversation_id);
CREATE INDEX idx_participants_user_id ON participants(user_id);
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_user_id ON messages(user_id);

-- View, make it easy to access the data
CREATE VIEW user_chats AS
WITH last_messages AS (
    SELECT
        m.conversation_id,
        m.id AS message_id,
        m.user_id AS message_user_id,
        m.type AS message_type,
        m.content AS message_content,
        m.created_at AS message_created_at
    FROM messages m
             INNER JOIN (
        SELECT
            conversation_id,
            MAX(created_at) AS last_message_time
        FROM messages
        GROUP BY conversation_id
    ) lm ON m.conversation_id = lm.conversation_id AND m.created_at = lm.last_message_time
)
SELECT
    c.id AS conversation_id,
    c.type AS conversation_type,
    c.name AS conversation_name,
    c.created_at AS conversation_created_at,
    c.updated_at AS conversation_updated_at,
    (
        SELECT JSON_GROUP_ARRAY(
                       JSON_OBJECT(
                               'id', u.id,
                               'display_name', u.display_name,
                               'email', u.email,
                               'role', p2.role
                       )
               )
        FROM users u
                 INNER JOIN participants p2 ON u.id = p2.user_id
        WHERE p2.conversation_id = c.id
    ) AS users,
    (
        SELECT JSON_GROUP_ARRAY(
                       JSON_OBJECT(
                               'id', lm.message_id,
                               'user_id', lm.message_user_id,
                               'type', lm.message_type,
                               'content', lm.message_content,
                               'created_at', lm.message_created_at
                       )
               )
        FROM last_messages lm
        WHERE lm.conversation_id = c.id
    ) AS messages
FROM
    conversations c
ORDER BY
    c.updated_at DESC;
