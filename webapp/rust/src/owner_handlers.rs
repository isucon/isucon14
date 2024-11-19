use std::collections::HashMap;

use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use axum_extra::extract::cookie::Cookie;
use axum_extra::extract::CookieJar;
use chrono::{DateTime, NaiveDate, Utc};

use crate::models::{Chair, Owner, Ride};
use crate::{AppState, Error};

pub fn owner_routes(app_state: AppState) -> axum::Router<AppState> {
    let routes =
        axum::Router::new().route("/api/owner/owners", axum::routing::post(owner_post_owners));

    let authed_routes = axum::Router::new()
        .route("/api/owner/sales", axum::routing::get(owner_get_sales))
        .route("/api/owner/chairs", axum::routing::get(owner_get_chairs))
        .route(
            "/api/owner/chairs/:chair_id",
            axum::routing::get(owner_get_chair_detail),
        )
        .route_layer(axum::middleware::from_fn_with_state(
            app_state.clone(),
            crate::middlewares::owner_auth_middleware,
        ));

    routes.merge(authed_routes)
}

#[derive(Debug, serde::Deserialize)]
struct OwnerPostOwnersRequest {
    name: String,
}

#[derive(Debug, serde::Serialize)]
struct OwnerPostOwnersResponse {
    id: String,
    chair_register_token: String,
}

async fn owner_post_owners(
    State(AppState { pool, .. }): State<AppState>,
    jar: CookieJar,
    axum::Json(req): axum::Json<OwnerPostOwnersRequest>,
) -> Result<(CookieJar, (StatusCode, axum::Json<OwnerPostOwnersResponse>)), Error> {
    let owner_id = ulid::Ulid::new().to_string();
    let access_token = crate::secure_random_str(32);
    let chair_register_token = crate::secure_random_str(32);

    sqlx::query(
        "INSERT INTO owners (id, name, access_token, chair_register_token) VALUES (?, ?, ?, ?)",
    )
    .bind(&owner_id)
    .bind(req.name)
    .bind(&access_token)
    .bind(&chair_register_token)
    .execute(&pool)
    .await?;

    let jar = jar.add(Cookie::build(("owner_session", access_token)).path("/"));

    Ok((
        jar,
        (
            StatusCode::CREATED,
            axum::Json(OwnerPostOwnersResponse {
                id: owner_id,
                chair_register_token,
            }),
        ),
    ))
}

#[derive(Debug, serde::Serialize)]
struct ChairSales {
    id: String,
    name: String,
    sales: i32,
}

#[derive(Debug, serde::Serialize)]
struct ModelSales {
    model: String,
    sales: i32,
}

#[derive(Debug, serde::Serialize)]
struct OwnerGetSalesResponse {
    total_sales: i32,
    chairs: Vec<ChairSales>,
    models: Vec<ModelSales>,
}

#[derive(Debug, serde::Deserialize)]
struct GetOwnerSalesQuery {
    since: Option<i64>,
    until: Option<i64>,
}

async fn owner_get_sales(
    State(AppState { pool, .. }): State<AppState>,
    axum::Extension(owner): axum::Extension<Owner>,
    Query(query): Query<GetOwnerSalesQuery>,
) -> Result<axum::Json<OwnerGetSalesResponse>, Error> {
    let since = if let Some(since) = query.since {
        DateTime::from_timestamp_millis(since).unwrap()
    } else {
        DateTime::from_timestamp_millis(0).unwrap()
    };
    let until = if let Some(until) = query.until {
        DateTime::from_timestamp_millis(until).unwrap()
    } else {
        DateTime::from_naive_utc_and_offset(
            NaiveDate::from_ymd_opt(9999, 12, 31)
                .unwrap()
                .and_hms_opt(23, 59, 59)
                .unwrap(),
            Utc,
        )
    };

    let mut tx = pool.begin().await?;

    let chairs: Vec<Chair> = sqlx::query_as("SELECT * FROM chairs WHERE owner_id = ?")
        .bind(&owner.id)
        .fetch_all(&mut *tx)
        .await?;

    let mut res = OwnerGetSalesResponse {
        total_sales: 0,
        chairs: Vec::with_capacity(chairs.len()),
        models: Vec::new(),
    };

    let mut model_sales_by_model = HashMap::new();

    for chair in chairs {
        let reqs: Vec<Ride> = sqlx::query_as("SELECT rides.* FROM rides JOIN ride_statuses ON rides.id = ride_statuses.ride_id WHERE chair_id = ? AND status = 'COMPLETED' AND updated_at BETWEEN ? AND ?")
            .bind(&chair.id)
            .bind(since)
            .bind(until)
            .fetch_all(&mut *tx)
            .await?;

        let chair_sales = sum_sales(&reqs);
        res.total_sales += chair_sales;

        res.chairs.push(ChairSales {
            id: chair.id,
            name: chair.name,
            sales: chair_sales,
        });

        *model_sales_by_model.entry(chair.model).or_insert(0) += chair_sales;
    }

    for (model, sales) in model_sales_by_model {
        res.models.push(ModelSales { model, sales });
    }

    Ok(axum::Json(res))
}

fn sum_sales(rides: &[Ride]) -> i32 {
    rides.iter().map(crate::calculate_sale).sum()
}

/// MySQL で COUNT()、SUM() 等を使って DECIMAL 型の値になったものを i64 に変換するための構造体。
#[derive(Debug)]
struct MysqlDecimal(i64);
impl sqlx::Decode<'_, sqlx::MySql> for MysqlDecimal {
    fn decode(
        value: sqlx::mysql::MySqlValueRef,
    ) -> Result<Self, Box<dyn std::error::Error + Send + Sync>> {
        use sqlx::{Type as _, ValueRef as _};

        let type_info = value.type_info();
        if i64::compatible(&type_info) {
            i64::decode(value).map(Self)
        } else if u64::compatible(&type_info) {
            let n = u64::decode(value)?.try_into()?;
            Ok(Self(n))
        } else if sqlx::types::Decimal::compatible(&type_info) {
            use num_traits::ToPrimitive as _;
            let n = sqlx::types::Decimal::decode(value)?
                .to_i64()
                .expect("failed to convert DECIMAL type to i64");
            Ok(Self(n))
        } else {
            panic!("MysqlDecimal is used with unknown type: {type_info:?}");
        }
    }
}
impl sqlx::Type<sqlx::MySql> for MysqlDecimal {
    fn type_info() -> sqlx::mysql::MySqlTypeInfo {
        i64::type_info()
    }

    fn compatible(ty: &sqlx::mysql::MySqlTypeInfo) -> bool {
        i64::compatible(ty) || u64::compatible(ty) || sqlx::types::Decimal::compatible(ty)
    }
}
impl From<MysqlDecimal> for i64 {
    fn from(value: MysqlDecimal) -> Self {
        value.0
    }
}

#[derive(Debug, sqlx::FromRow)]
struct ChairWithDetail {
    id: String,
    owner_id: String,
    name: String,
    access_token: String,
    model: String,
    is_active: bool,
    created_at: DateTime<Utc>,
    updated_at: DateTime<Utc>,
    total_distance: MysqlDecimal,
    total_distance_updated_at: Option<DateTime<Utc>>,
}

#[derive(Debug, serde::Serialize)]
struct OwnerChair {
    id: String,
    name: String,
    model: String,
    active: bool,
    registered_at: i64,
    total_distance: i64,
    #[serde(skip_serializing_if = "Option::is_none")]
    total_distance_updated_at: Option<i64>,
}

#[derive(Debug, serde::Serialize)]
struct OwnerGetChairResponse {
    chairs: Vec<OwnerChair>,
}

async fn owner_get_chairs(
    State(AppState { pool, .. }): State<AppState>,
    axum::Extension(owner): axum::Extension<Owner>,
) -> Result<axum::Json<OwnerGetChairResponse>, Error> {
    let chairs: Vec<ChairWithDetail> = sqlx::query_as(r#"SELECT id,
       owner_id,
       name,
       access_token,
       model,
       is_active,
       created_at,
       updated_at,
       IFNULL(total_distance, 0) AS total_distance,
       total_distance_updated_at
FROM chairs
       LEFT JOIN (SELECT chair_id,
                          SUM(IFNULL(distance, 0)) AS total_distance,
                          MAX(created_at)          AS total_distance_updated_at
                   FROM (SELECT chair_id,
                                created_at,
                                ABS(latitude - LAG(latitude) OVER (PARTITION BY chair_id ORDER BY created_at)) +
                                ABS(longitude - LAG(longitude) OVER (PARTITION BY chair_id ORDER BY created_at)) AS distance
                         FROM chair_locations) tmp
                   GROUP BY chair_id) distance_table ON distance_table.chair_id = chairs.id
WHERE owner_id = ?
    "#).bind(owner.id).fetch_all(&pool).await?;

    Ok(axum::Json(OwnerGetChairResponse {
        chairs: chairs
            .into_iter()
            .map(|chair| OwnerChair {
                id: chair.id,
                name: chair.name,
                model: chair.model,
                active: chair.is_active,
                registered_at: chair.created_at.timestamp_millis(),
                total_distance: chair.total_distance.0,
                total_distance_updated_at: chair
                    .total_distance_updated_at
                    .map(|t| t.timestamp_millis()),
            })
            .collect(),
    }))
}

#[derive(Debug, serde::Serialize)]
struct OwnerGetChairDetailResponse {
    id: String,
    name: String,
    model: String,
    active: bool,
    registered_at: i64,
    total_distance: i64,
    total_distance_updated_at: Option<i64>,
}

async fn owner_get_chair_detail(
    State(AppState { pool, .. }): State<AppState>,
    Path((chair_id,)): Path<(String,)>,
    axum::Extension(owner): axum::Extension<Owner>,
) -> Result<axum::Json<OwnerGetChairDetailResponse>, Error> {
    let Some(chair): Option<ChairWithDetail> = sqlx::query_as(r#"SELECT id,
       owner_id,
       name,
       access_token,
       model,
       is_active,
       created_at,
       updated_at,
       IFNULL(total_distance, 0) AS total_distance,
       total_distance_updated_at
FROM chairs
       LEFT JOIN (SELECT chair_id,
                          SUM(IFNULL(distance, 0)) AS total_distance,
                          MAX(created_at)          AS total_distance_updated_at
                   FROM (SELECT chair_id,
                                created_at,
                                ABS(latitude - LAG(latitude) OVER (PARTITION BY chair_id ORDER BY created_at)) +
                                ABS(longitude - LAG(longitude) OVER (PARTITION BY chair_id ORDER BY created_at)) AS distance
                         FROM chair_locations) tmp
                   GROUP BY chair_id) distance_table ON distance_table.chair_id = chairs.id
WHERE owner_id = ? AND id = ?"#)
    .bind(owner.id)
    .bind(chair_id)
    .fetch_optional(&pool)
    .await? else {
        return Err(Error::NotFound("chair not found"));
    };

    Ok(axum::Json(OwnerGetChairDetailResponse {
        id: chair.id,
        name: chair.name,
        model: chair.model,
        active: chair.is_active,
        registered_at: chair.created_at.timestamp_millis(),
        total_distance: chair.total_distance.0,
        total_distance_updated_at: chair
            .total_distance_updated_at
            .map(|t| t.timestamp_millis()),
    }))
}
