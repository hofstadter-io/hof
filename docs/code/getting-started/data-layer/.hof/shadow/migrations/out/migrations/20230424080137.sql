/* has snapshot 0 - 20230424080137*/
/* == User == */
CREATE TABLE
  user (
    created_at datetime,
    id uuid DEFAULT uuid_generate_v4 (),
    profile uuid,
    updated_at datetime,
    active bool DEFAULT false,
    email varchar(64),
    password varchar(64),
    username varchar(64),
  );

/* == UserProfile == */
CREATE TABLE
  user_profile (
    about varchar(64),
    avatar varchar(64),
    created_at datetime,
    id uuid DEFAULT uuid_generate_v4 (),
    owner uuid,
    updated_at datetime,
  );
