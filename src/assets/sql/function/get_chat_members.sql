CREATE OR REPLACE FUNCTION public.get_chat_members(p_chat_id bigint)
    RETURNS SETOF chats
    LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    SELECT 
        c.id,
        CASE 
            WHEN c.type = 'private' THEN u.login
            ELSE c.name
        END,
        CASE 
            WHEN c.type = 'private' THEN u.avatar_url
            ELSE c.avatar_url
        END,
        c.type
    FROM chats c
    INNER JOIN chat_members cm ON c.id = cm.chat_id
    LEFT JOIN LATERAL (
        SELECT cm2.user_id
        FROM chat_members cm2
        WHERE cm2.chat_id = c.id 
          AND cm2.user_id != p_user_id
        LIMIT 1
    ) friend ON true
    LEFT JOIN users u ON friend.user_id = u.id
    WHERE cm.user_id = p_user_id;
END;
$function$
