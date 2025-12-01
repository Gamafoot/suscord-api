CREATE OR REPLACE FUNCTION tg_remove_friend()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.user_id != OLD.friend_id AND EXISTS (
        SELECT 1 FROM friends 
        WHERE user_id = OLD.friend_id AND friend_id = OLD.user_id
    ) THEN
        DELETE FROM friends WHERE user_id = OLD.friend_id AND friend_id = OLD.user_id;
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tg_remove_friend ON friends;

CREATE TRIGGER tg_remove_friend
AFTER DELETE
ON friends
FOR EACH ROW
EXECUTE PROCEDURE tg_remove_friend();
