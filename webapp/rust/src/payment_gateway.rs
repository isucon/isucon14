use crate::models::Ride;
use crate::Error;
use std::future::Future;

#[derive(Debug, thiserror::Error)]
pub enum PaymentGatewayError {
    #[error("reqwest error: {0}")]
    Reqwest(#[from] reqwest::Error),
    #[error("unexpected number of payments: {ride_count} != {payment_count}.")]
    UnexpectedNumberOfPayments {
        ride_count: usize,
        payment_count: usize,
    },
    #[error("[GET /payments] unexpected status code ({0})")]
    GetPayment(reqwest::StatusCode),
}

#[derive(Debug, serde::Serialize)]
pub struct PaymentGatewayPostPaymentRequest {
    pub amount: i32,
}

#[derive(Debug, serde::Deserialize)]
struct PaymentGatewayGetPaymentsResponseOne {
    amount: i32,
    status: String,
}

pub trait PostPaymentCallback<'a> {
    type Output: Future<Output = Result<Vec<Ride>, Error>>;

    fn call(&self, tx: &'a mut sqlx::MySqlConnection, user_id: &'a str) -> Self::Output;
}
impl<'a, F, Fut> PostPaymentCallback<'a> for F
where
    F: Fn(&'a mut sqlx::MySqlConnection, &'a str) -> Fut,
    Fut: Future<Output = Result<Vec<Ride>, Error>>,
{
    type Output = Fut;
    fn call(&self, tx: &'a mut sqlx::MySqlConnection, user_id: &'a str) -> Fut {
        self(tx, user_id)
    }
}

pub async fn request_payment_gateway_post_payment<F>(
    payment_gateway_url: &str,
    token: &str,
    param: &PaymentGatewayPostPaymentRequest,
    tx: &mut sqlx::MySqlConnection,
    user_id: &str,
    retrieve_rides_order_by_created_at_asc: F,
) -> Result<(), Error>
where
    F: for<'a> PostPaymentCallback<'a>,
{
    // 失敗したらとりあえずリトライ
    // FIXME: 社内決済マイクロサービスのインフラに異常が発生していて、同時にたくさんリクエストすると変なことになる可能性あり
    let mut retry = 0;

    loop {
        let result = async {
            let res = reqwest::Client::new()
                .post(format!("{payment_gateway_url}/payments"))
                .bearer_auth(token)
                .json(param)
                .send()
                .await
                .map_err(PaymentGatewayError::Reqwest)?;

            if res.status() != reqwest::StatusCode::NO_CONTENT {
                // エラーが返ってきても成功している場合があるので、社内決済マイクロサービスに問い合わせ
                let get_res = reqwest::Client::new()
                    .get(format!("{payment_gateway_url}/payments"))
                    .bearer_auth(token)
                    .send()
                    .await
                    .map_err(PaymentGatewayError::Reqwest)?;

                // GET /payments は障害と関係なく200が返るので、200以外は回復不能なエラーとする
                if get_res.status() != reqwest::StatusCode::OK {
                    return Err(PaymentGatewayError::GetPayment(get_res.status()).into());
                }
                let payments: Vec<PaymentGatewayGetPaymentsResponseOne> =
                    get_res.json().await.map_err(PaymentGatewayError::Reqwest)?;

                let rides = retrieve_rides_order_by_created_at_asc
                    .call(tx, user_id)
                    .await?;

                if rides.len() != payments.len() {
                    return Err(PaymentGatewayError::UnexpectedNumberOfPayments {
                        ride_count: rides.len(),
                        payment_count: payments.len(),
                    }
                    .into());
                }
            }
            Ok(())
        }
        .await;

        if let Err(err) = result {
            if retry < 5 {
                retry += 1;
                tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
                continue;
            } else {
                return Err(err);
            }
        }
        break;
    }

    Ok(())
}
