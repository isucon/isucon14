use axum::extract::{Request, State};
use axum::middleware::Next;
use axum::response::Response;
use axum_extra::extract::CookieJar;

use crate::models::{Chair, Owner, User};
use crate::{AppState, Error};

pub async fn app_auth_middleware(
    State(AppState { pool, .. }): State<AppState>,
    jar: CookieJar,
    mut req: Request,
    next: Next,
) -> Result<Response, Error> {
    let Some(c) = jar.get("app_session") else {
        return Err(Error::Unauthorized("app_session cookie is required"));
    };
    let access_token = c.value();
    let Some(user): Option<User> = sqlx::query_as("SELECT * FROM users WHERE access_token = ?")
        .bind(access_token)
        .fetch_optional(&pool)
        .await?
    else {
        return Err(Error::BadRequest("invalid access token"));
    };

    req.extensions_mut().insert(user);

    Ok(next.run(req).await)
}

pub async fn owner_auth_middleware(
    State(AppState { pool, .. }): State<AppState>,
    jar: CookieJar,
    mut req: Request,
    next: Next,
) -> Result<Response, Error> {
    let Some(c) = jar.get("owner_session") else {
        return Err(Error::Unauthorized("owner_session cookie is required"));
    };
    let access_token = c.value();
    let Some(owner): Option<Owner> = sqlx::query_as("SELECT * FROM owners WHERE access_token = ?")
        .bind(access_token)
        .fetch_optional(&pool)
        .await?
    else {
        return Err(Error::BadRequest("invalid access token"));
    };

    req.extensions_mut().insert(owner);

    Ok(next.run(req).await)
}

pub async fn chair_auth_middleware(
    State(AppState { pool, .. }): State<AppState>,
    jar: CookieJar,
    mut req: Request,
    next: Next,
) -> Result<Response, Error> {
    let Some(c) = jar.get("chair_session") else {
        return Err(Error::Unauthorized("chair_session cookie is required"));
    };
    let access_token = c.value();
    let Some(chair): Option<Chair> = sqlx::query_as("SELECT * FROM chairs WHERE access_token = ?")
        .bind(access_token)
        .fetch_optional(&pool)
        .await?
    else {
        return Err(Error::BadRequest("invalid access token"));
    };

    req.extensions_mut().insert(chair);

    Ok(next.run(req).await)
}
