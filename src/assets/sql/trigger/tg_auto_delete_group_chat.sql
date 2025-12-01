CREATE OR REPLACE FUNCTION tg_auto_delete_group_chat()
RETURNS TRIGGER AS $$
DECLARE
    members_count int;
BEGIN
    SELECT COUNT(*) INTO members_count FROM chat_members WHERE chat_id = OLD.chat_id;
    IF members_count = 0 THEN
        DELETE FROM chats WHERE id = OLD.chat_id;
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tg_auto_delete_group_chat ON chat_members;

CREATE TRIGGER tg_auto_delete_group_chat
AFTER DELETE
ON chat_members
FOR EACH ROW
EXECUTE PROCEDURE tg_auto_delete_group_chat();
