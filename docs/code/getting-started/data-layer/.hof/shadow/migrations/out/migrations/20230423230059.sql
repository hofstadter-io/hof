/* has snapshot 0 - 20230423230059*/
/* first datamodel, first snapshot */
CREATE TABLE
  user (
    created_at datetime,
    deleted_at datetime,
    id uuid,
    profile uuid,
    updated_at datetime,
    active bool,
    email varchar(64),
    password varchar(64),
    username varchar(64),
  );

CREATE TABLE
  user_profile (
    about varchar(64),
    avatar varchar(64),
    owner uuid,
  );
