CREATE OR REPLACE FUNCTION public.get_user_chats(p_chat_id bigint)
    RETURNS SETOF users
    LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY SELECT * FROM users WHERE id IN (
        SELECT user_id FROM chat_members 
        WHERE chat_id = p_chat_id
    );
END;
$function$
