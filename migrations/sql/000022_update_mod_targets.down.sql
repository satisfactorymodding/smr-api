Create or replace function update_mod_platform_down_random_string(length integer) returns text as
$$
declare
    chars text[] := '{0,1,2,3,4,5,6,7,8,9,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z}';
    result text := '';
    i integer := 0;
begin
    if length < 0 then
        raise exception 'Given length cannot be less than 0';
    end if;
    for i in 1..length loop
            result := result || chars[1+random()*(array_length(chars, 1)-1)];
        end loop;
    return result;
end;
$$ language plpgsql;

ALTER TABLE version_targets RENAME TO mod_archs;

ALTER TABLE mod_archs
    RENAME COLUMN version_id TO mod_version_arch_id;
ALTER TABLE mod_archs
    RENAME COLUMN target_name TO platform;
ALTER TABLE mod_archs
    ADD COLUMN id varchar(14);

-- This is not as random as the original ID, but it should be good enough
UPDATE mod_archs SET id = update_mod_platform_down_random_string(14) WHERE true;

ALTER TABLE mod_archs
    ALTER COLUMN id SET NOT NULL;

ALTER TABLE mod_archs
    DROP CONSTRAINT version_targets_version_id_fkey,
    DROP CONSTRAINT version_targets_pkey,
    ADD CONSTRAINT mod_archs_pkey PRIMARY KEY (id);

CREATE INDEX IF NOT EXISTS idx_mod_arch_id ON mod_archs (mod_version_arch_id, platform);

ALTER TABLE sml_version_targets RENAME TO sml_archs;

ALTER TABLE sml_archs
    RENAME COLUMN version_id TO sml_version_arch_id;
ALTER TABLE sml_archs
    RENAME COLUMN target_name TO platform;
ALTER TABLE sml_archs
    ADD COLUMN id varchar(14);

-- This is not as random as the original ID, but it should be good enough
UPDATE sml_archs SET id = update_mod_platform_down_random_string(14) WHERE true;

ALTER TABLE sml_archs
    ALTER COLUMN id SET NOT NULL;

ALTER TABLE sml_archs
    DROP CONSTRAINT sml_version_targets_version_id_fkey,
    DROP CONSTRAINT sml_version_targets_pkey,
    ADD CONSTRAINT sml_archs_pkey PRIMARY KEY (id);

CREATE INDEX IF NOT EXISTS idx_sml_archs_id ON sml_archs (sml_version_arch_id, platform);

DROP FUNCTION update_mod_platform_down_random_string(length integer);