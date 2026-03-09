-- name: ListHumidityLog :many
select *
from humidity_log
order by time desc
limit ? offset ?;

-- name: HumidityLogForDeviceToDate :many
select *
from humidity_log
         join device_tag on device_name = humidity_log.sensor
where device_tag.tag = ?
  and (time between ? and ?)
order by time desc;


-- name: GetHumidityLog :one
select *
from humidity_log
where id = ?
limit 1;

-- name: CreateHumidityLog :one
insert
into humidity_log (time, sensor, temperature, humidity, dew_point)
values (?, ?, ?, ?, ?)
returning *;

-- name: DeleteHumidityLog :exec
delete
from humidity_log
where id = ?;


-- name: ListPowerLog :many
select *
from power_log
order by time desc
limit ? offset ?;

-- name: PowerLogForDeviceToDate :many
select *
from power_log
         join device_tag on device_name = power_log.sensor
where device_tag.tag = ?
  and (time between ? and ?)
order by time desc;


-- name: GetPowerLog :one
select *
from power_log
where id = ?
limit 1;

-- name: CreatePowerLog :one
insert
into power_log (time, sensor, total_start_time, total, yesterday, today, period, power, apparent_power, reactive_power,
                factor, voltage, current, sensor_temperature)
values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
returning *;

-- name: DeletePowerLog :exec
delete
from power_log
where id = ?;

-- name: ListSwitchStateLog :many
select *
from switch_state_log
order by time desc
limit ? offset ?;

-- name: SwitchStateLogForDeviceToDate :many
select *
from switch_state_log
         join device_tag on device_name = switch_state_log.sensor
where device_tag.tag = ?
  and (time between ? and ?)
order by time desc;

-- name: GetSwitchStateLog :one
select *
from switch_state_log
where id = ?
limit 1;

-- name: CreateSwitchStateLog :one
insert
into switch_state_log (time, sensor, switch_state)
values (?, ?, ?)
returning *;

-- name: DeleteSwitchStateLog :exec
delete
from switch_state_log
where id = ?;

-- name: ListDeviceTag :many
select *
from device_tag;

-- name: ListDevices :many
select tag
from device_tag
group by tag;

-- name: GetDeviceTag :one
select *
from device_tag
where id = ?
limit 1;

-- name: CreateDeviceTag :one
insert
into device_tag (device_name, tag)
values (?, ?)
returning *;

-- name: UpdateDeviceTag :one
UPDATE device_tag
set device_name = ?,
    tag         = ?
WHERE id = ?
RETURNING *;

-- name: DeleteDeviceTag :exec
delete
from device_tag
where id = ?;

-- name: TruncateDeviceTag :exec
delete
from device_tag;



