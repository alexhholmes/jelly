create table users
(
    id uuid default gen_random_uuid() not null,
    username VARCHAR(50) UNIQUE,
    email VARCHAR(255) UNIQUE,
    profile_image_url TEXT,
    bio TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    constraint users_pk
        primary key (id)
);

create table raw_photos
(
    id                uuid default gen_random_uuid()         not null,
    user_id           uuid                                   not null,
    original_filename varchar(255)                           not null,
    storage_url       varchar(500)                           not null,
    file_size         bigint                                 not null,
    mime_type         varchar(100)                           not null,
    md5_hash          varchar(32)                            not null,
    width             integer,
    height            integer,
    exif_data         jsonb,
    uploaded_at       timestamp with time zone default now() not null,
    processed_at      timestamp with time zone,
    schedule_deletion timestamp with time zone,
    constraint raw_photos_pk
        primary key (id),
    constraint raw_photos_user_fk
        foreign key (user_id) references users (id),
    constraint raw_photos_md5_unique
        unique (md5_hash)
);

create table photos
(
    id                uuid default gen_random_uuid()         not null,
    raw_photo_id      uuid                                   not null,
    user_id           uuid                                   not null,
    filename          varchar(255)                           not null,
    original_url      varchar(500)                           not null,
    thumbnail_url     varchar(500)                           not null,
    caption           text,
    tags              text[],
    file_size         bigint                                 not null,
    mime_type         varchar(100)                           not null,
    width             integer,
    height            integer,
    uploaded_at       timestamp with time zone default now() not null,
    updated_at        timestamp with time zone default now() not null,
    schedule_deletion timestamp with time zone,
    constraint photos_pk
        primary key (id),
    constraint photos_raw_photo_fk
        foreign key (raw_photo_id) references raw_photos (id),
    constraint photos_user_fk
        foreign key (user_id) references users (id)
);

create table photo_likes
(
    photo_id    UUID REFERENCES photos(id),
    user_id     UUID REFERENCES users(id),
    created_at  TIMESTAMP,
    PRIMARY KEY (photo_id, user_id)
);

create table photo_comments
(
    id UUID PRIMARY KEY,
    photo_id UUID REFERENCES photos(id),
    user_id UUID REFERENCES users(id),
    content TEXT,
    created_at TIMESTAMP
);