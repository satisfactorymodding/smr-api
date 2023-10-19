--Mod Targets--
DELETE FROM version_targets
    USING versions
    WHERE version_targets.version_id = versions.id AND
            version_targets.target_name = 'Windows' AND
            version_targets.key = versions.key AND
            version_targets.hash = versions.hash AND
            version_targets.size = versions.size;

--SML Targets--
DELETE FROM sml_version_targets
    USING sml_versions
    WHERE sml_version_targets.version_id = sml_versions.id AND
            sml_version_targets.target_name = 'Windows' AND
            sml_version_targets.link = replace(sml_versions.link, '/tag/', '/download/') || '/SML.zip';
