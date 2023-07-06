--Mod Targets--
DELETE FROM version_targets
WHERE EXISTS ( 
    SELECT 1 
    FROM versions 
    WHERE 
        version_targets.version_id = versions.id AND 
        version_targets.platform = 'Windows' AND 
        version_targets.key = versions.key AND 
        version_targets.hash = versions.hash AND 
        version_targets.size = versions.size
    );

--SML Targets--
DELETE FROM sml_version_targets
WHERE EXISTS ( 
    SELECT 1 
    FROM sml_versions 
    WHERE 
        sml_version_targets.version_id = sml_versions.id AND 
        sml_version_targets.platform = 'Windows' AND 
        sml_version_targets.link = sml_versions.link
    );