UPDATE "mods"
  SET "compatibility" = mods.compatibility - 'Controller'
  WHERE mods.compatibility ? 'Controller' AND mods.compatibility ? 'EA';
