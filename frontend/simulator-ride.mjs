const chair = {
  "id": "01JDFEF7MGXXCJKW1MNJXPA77A",
  "owner_id": "01JDFEDF008NTA922W12FS7800",
  "name": "QC-L13-8361",
  "model": "クエストチェア Lite",
  "token": "3013d5ec84e1b230f913a17d71ef27c8"
};

const user = {
  "id": "01JDJ23EA0C0P2KFPTXDKTZMNM",
  "username": "Collier6283",
  "firstname": "和治",
  "lastname": "大森",
  "date_of_birth": "1972-01-20",
  "token": "34ea320039fc61ae2558176607a2e12c",
  "invitation_code": "5c4a695f66d598e"
}

const apiURL = "http://localhost:8080"

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}


async function postRide () {
  const rides = await fetch(`${apiURL}/api/app/rides`, {
    "method": "POST",
    "headers": {
      "Cookie": `app_session=${user.token};`
    },
    "body": JSON.stringify({"pickup_coordinate":{"latitude":-4,"longitude":2},"destination_coordinate":{"latitude":2,"longitude":10}})
  })
  return rides
}

async function getAppNotification() {
   const appNotification = await fetch(`${apiURL}/api/app/notification`, {
    "method": "GET",
    "headers": {
      "Cookie": `app_session=${user.token};`
    }
  })
  return appNotification.json()
}

async function getChairNotification() {
  try {
  console.log("chairNotification")
  const chairNotification = await fetch(`${apiURL}/api/chair/notification`, {
   "method": "GET",
   "headers": {
     "Cookie": `chair_session=${chair.token};`
   }
 })
 console.log("chairNotificationaaa", chairNotification)
 return await chairNotification.json()
  } catch(e) {
    console.log("error")
    console.error(e)
  }
}

async function chairCooridnate ({latitude, longitude}) {

const coordinate = await fetch(`${apiURL}/api/chair/coordinate`, {
  "method": "POST",
  "headers": {
    "Cookie": `chair_session=${chair.token}`
  },
  "body": JSON.stringify({"latitude": latitude,"longitude":longitude})
})
  return coordinate
}
async function chairRides(rideId, setStatus) {
  const status = await fetch(`${apiURL}/api/chair/rides/${rideId}/status`, {
    "method": "POST",
    "headers": {
      "Cookie": `chair_session=${chair.token}`
    },
    "body": JSON.stringify({
      "status": setStatus
    })
  })
  return status
}

const main = async() => {
  /**
 * activate
 */
const activate = await fetch(`${apiURL}/api/chair/activity`, {
  "method": "POST",
  "headers": {
    "Cookie": `chair_session=${chair.token};`
  },
  "body": JSON.stringify({"is_active":true})
})

console.log({activate})

/**
 * coordinate
 */
const coordinate = await chairCooridnate({latitude: 1,longitude: 20})

console.log({coordinate})

/**
 * rideリクエスト
 */
try {
const requested = await postRide()

console.log({requested})
} catch (e) {
  console.error(e)
}

await sleep(2000)

const notification = await getChairNotification()

console.log(notification)
const rideId = notification.data.ride_id

/**
 * ENROUTE
 */
const status = await chairRides(rideId, "ENROUTE")

console.log({status})

await sleep(2000)

const notification2 = (await getChairNotification()).data

const pickup = notification2.pickup_coordinate

console.log({pickup})

await chairCooridnate(pickup)
await sleep(2000)

await chairRides(rideId, "CARRYING")
await sleep(2000)

const notification3 = (await getChairNotification()).data
const destination_coordinate = notification3.destination_coordinate
await chairCooridnate(destination_coordinate)

await sleep(30000)
console.log('sleep')
/**
 * evaluate
 */
try {
  const ride_id = (await getAppNotification()).data.ride_id
  fetch(`${apiURL}/api/app/rides/${ride_id}/evaluation`, {
    "method": "POST",
    "headers": {
      "Cookie": `app_session=${user.token};`
    },
    "body": JSON.stringify({"evaluation":4})
  })
  postRide()
  } catch (e) {
    console.error(e)
  }
  await sleep(30000)
  const clientNotifiaction = (await getAppNotification())
  const chairNotification = await getChairNotification()

  console.log({clientNotifiaction, chairNotification})
}

main()