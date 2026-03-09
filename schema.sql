create table if not exists humidity_log
(
    id          integer primary key,
    time        datetime not null,
    sensor      text     not null,
    temperature real     not null,
    humidity    real     not null,
    dew_point   real     not null
);

create table if not exists power_log
(
    id                 integer primary key,
    time               datetime not null,
    sensor             text     not null,
    total_start_time   datetime not null,
    total              real     not null,
    yesterday          real     not null,
    today              real     not null,
    period             real     not null,
    power              real     not null,
    apparent_power     real     not null,
    reactive_power     real     not null,
    factor             real     not null,
    voltage            real     not null,
    current            real     not null,
    sensor_temperature real     not null
);

create table if not exists switch_state_log
(
    id                 integer primary key,
    time               datetime not null,
    sensor             text     not null,
    switch_state       text not null
);

create table if not exists device_tag
(
    id integer primary key,
    device_name text not null,
    tag text not null
)