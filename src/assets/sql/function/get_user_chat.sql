CREATE OR REPLACE FUNCTION public.get_user_chat(p_chat_id bigint, p_user_id bigint)
 RETURNS TABLE(id bigint, name character varying, avatar_url character varying, type character varying)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY 
    SELECT * FROM get_user_chats(p_user_id) AS chats
    WHERE chats.id = p_chat_id;
END;
$function$
