CREATE OR REPLACE FUNCTION public.get_user_chats(p_user_id bigint)
 RETURNS SETOF chat_record
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    SELECT 
        c.id,
        CASE WHEN c.type = 'private' THEN u.username ELSE c.name END AS name,
        CASE WHEN c.type = 'private' THEN u.avatar_path ELSE c.avatar_path END AS avatar_path,
        c.type
    FROM chats c
    JOIN chat_members cm ON cm.chat_id = c.id
    LEFT JOIN LATERAL (
        SELECT cm2.user_id
        FROM chat_members cm2
        WHERE cm2.chat_id = c.id AND cm2.user_id != p_user_id
        LIMIT 1
    ) friend ON true
    LEFT JOIN users u ON friend.user_id = u.id
    LEFT JOIN LATERAL (
        SELECT created_at 
        FROM messages m 
        WHERE m.chat_id = c.id AND m.user_id = p_user_id
        ORDER BY m.created_at DESC
        LIMIT 1
    ) last_msg ON true
    WHERE cm.user_id = p_user_id
    ORDER BY last_msg.created_at DESC NULLS LAST;
END;
$function$
