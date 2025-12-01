CREATE OR REPLACE FUNCTION tg_add_friend()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.user_id != NEW.friend_id AND NOT EXISTS (
        SELECT 1 FROM friends 
        WHERE user_id = NEW.friend_id AND friend_id = NEW.user_id
    ) THEN
        INSERT INTO friends (user_id, friend_id) VALUES (NEW.friend_id, NEW.user_id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tg_add_friend ON friends;

CREATE TRIGGER tg_add_friend
AFTER INSERT
ON friends
FOR EACH ROW
EXECUTE PROCEDURE tg_add_friend();
