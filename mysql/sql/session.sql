#namespace session_db

create database session charset utf8;

#create tables
/*
create table if not exists t_session_info (
    sessionuuid varchar(64),
    timeoutS bigint,
    losevalidtime bigint
) charset utf8;
*/
#end

/*
    @bref 添加session
    @is_brace true
    @in_isarr false
    @out_isarr false
    @in sessionUuid: string
    @in timeoutS: int64
    @in loseValidTime: int64
*/
#define addSession
insert into t_session_info values({0}, {1}, {2});
#end

/*
    @bref 删除session
    @is_brace true
    @in_isarr true
    @out_isarr false
    @in sessionUuid: string
*/
#define deleteSession
delete from t_session_info where sessionuuid = {0};
#end

/*
    @bref 更新session
    @is_brace true
    @in_isarr true
    @out_isarr false
    @in condition[cond]: string
    @in sessionUuid: string
*/
#define updateSession
update t_session_info set {0} where sessionuuid = {1};
#end

/*
    @bref 获取单个session
    @is_brace true
    @in_isarr false
    @out_isarr false
    @in sessionUuid: string
    @out timeoutS: int64
    @out loseValidTime: int64
*/
#define getSession
select timeoutS, losevalidtime from t_session_info
where sessionuuid = {0};
#end

/*
    @bref 根据session获取个数
    @is_brace true
    @in_isarr false
    @out_isarr false
    @in sessionUuid: string
    @out count: int
*/
#define getCountBySessionUuid
select count(0) from t_session_info
where sessionuuid = {0};
#end

/*
    @bref 删除过期时间记录
    @is_brace true
    @in_isarr false
    @out_isarr false
*/
#define deleteLosetimeRecord
delete from t_session_info where sessionuuid in
(
    select tmp.sessionuuid from
    (
        select sessionuuid from t_session_info where losevalidtime < unix_timestamp(now())
    ) as tmp
);
#end
