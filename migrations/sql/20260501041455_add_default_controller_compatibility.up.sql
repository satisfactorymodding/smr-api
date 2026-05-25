UPDATE "mods"
  -- Assign Controller compatibility of Untested...
  SET "compatibility" = jsonb_set(mods.compatibility, '{Controller}', '{"Note":"","State":"Untested"}'::jsonb)
  -- ... to mods with pre-controller Compatibility already set + no leftover existing controller compatibility info
  WHERE compatibility ? 'EA' and (NOT compatibility ? 'Controller');
