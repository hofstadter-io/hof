/* has snapshot 2 - 20230424080345*/
/* == User == */
ALTER TABLE
  user ADD posts uuid,;

/* == UserPost == */
CREATE TABLE
  user_post (
    content varchar(2048),
    created_at datetime,
    format varchar(64),
    id uuid DEFAULT uuid_generate_v4 (),
    owner uuid,
    title varchar(64),
    updated_at datetime,
  );

/* == UserProfile == */
ALTER TABLE
  user_profile ADD social varchar(64),;
