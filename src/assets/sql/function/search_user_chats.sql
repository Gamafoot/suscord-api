CREATE OR REPLACE FUNCTION public.search_user_chats(p_user_id bigint, p_search_pattern text DEFAULT NULL::text)
 RETURNS SETOF chat_record
 LANGUAGE plpgsql
AS $function$
BEGIN 
    RETURN QUERY
    SELECT id, name, avatar_path, type FROM get_user_chats(p_user_id) WHERE (p_search_pattern IS NULL OR name ~* p_search_pattern);
END;
$function$
