/**
 * 初期データが作成されるまで一旦こちらで作成してdata.jsonに突っ込むようにする
 */

import { writeFileSync } from "fs";

const create = async() => {
  const candidates = [
    {name: "オーナー1", chairs: [{
      name: "chair1", model: "model1"
    }, {
      name: "chair2", model: "model2"
    }, {
      name: "chair3", model: "model3"
    }]},
    {name: "オーナー2", chairs: [
      {
        name: "chair1", model: "model1"
      }
    ]},
    {name: "オーナー3", chairs: []}
  ]

  const BASE_URL = "http://localhost:8080";

  return candidates.map(async (candidate) => {
    const ownerFetch = await fetch(`${BASE_URL}/api/owner/owners`, {
      "body": JSON.stringify({
        name: candidate.name
      }),
      "credentials": 'include',
      "method": "POST"
    })
    /**
     * @type {{id: string,chair_register_token: string}}
     */
    const json = await ownerFetch.json();

    const chairRegisterToken = json.chair_register_token;

    console.log('check: ownerFetch', ownerFetch.headers["Cookie"])
    // set-cookie ヘッダーを取得
    const cookies = ownerFetch.headers.raw()['set-cookie'];
    console.log('All Cookies:', cookies);

    // owner_session クッキーを探す
    const ownerSessionCookie = cookies.find(cookie => cookie.startsWith('owner_session='));
    console.log('owner_session Cookie:', ownerSessionCookie);
    const ownerSessionValue = ownerSessionCookie
      ? ownerSessionCookie.split(';')[0].split('=')[1]
      : null;
    
    const chairs = candidate.chairs.map( async (chair) => {
      const chairFetch = await fetch(`${BASE_URL}/api/owner/chairs`, {
        body: JSON.stringify({name: chair.name, model: chair.model, chair_register_token: chairRegisterToken}),
        method: "POST",
        "credentials": 'include',
      })
      /**
       * @type {{id: string, owner_id: string}}
       */
      const json = await chairFetch.json()

      // set-cookie ヘッダーを取得
      console.log('headers', chairFetch.headers)
      const cookies =  chairFetch.headers.raw()['set-cookie'];
      console.log('All Cookies:', cookies);

      // owner_session クッキーを探す
      const chairSessionCookie = cookies.find(cookie => cookie.startsWith('chair_session='));
      console.log('chair_session Cookie:', chairSessionCookie);
      const chairSessionValue = chairSessionCookie
        ? chairSessionCookie.split(';')[0].split('=')[1]
        : null;
      return {
        id: json.id,
        name: chair.name,
        model: chair.model,
        token: chairSessionValue
      }
    })

    return {
      id: json.id,
      name: candidate.name,
      token: ownerSessionValue,
      chairs
    }
  })

}


const main = async() => {
  const data = await create()
  console.log('data', data)
  writeFileSync("./app/initial-data/data.json", JSON.stringify(data, null, 2))
}


main()