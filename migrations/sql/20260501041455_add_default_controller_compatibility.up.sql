UPDATE "mods"
  SET "compatibility" = mods.compatibility || '{"Controller":""}'
  WHERE NOT mods.compatibility ? 'Controller' AND mods.compatibility ? 'EA';

UPDATE "mods"
  SET "compatibility" = jsonb_set(mods.compatibility, '{Controller}', '{"Note":"","State":"Untested"}'::jsonb)
  WHERE mods.compatibility @> '{"Controller":""}'